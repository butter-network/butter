package butter

import (
	"errors"
	"fmt"
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/utils"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pbnjay/memory"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	// Listener communication byte codes
	AppCode   byte = 100 // received a request to interact with the app server behaviour
	pingCode  byte = 101 // received a ping request from a node in startup mode
	helloCode byte = 102 // received a hello request from a node in response to a ping
	success   byte = 200 // received a success response from a node
	failure   byte = 99  // received a failure response from a node
	GroupCode byte = 20  // received a request to get the groups of a node
)
const BlockSize = 4096

// A Block has a size of 4096 bytes with a uuid of size 16 bytes, 5 keywords of max size 50 bytes, 2 part numbers of
// size 8 bytes each, a geotag of size 2 bytes and a data of size 4096 - 16 - 5*50 - 2*8 - 2 = 3840 bytes. A Block is
// uniquely identified by combining its uuid and part number e.g. <UUID>/<PartNumber>.
type Block struct {
	uuid     [16]byte    // probably don't need this?
	keywords [5][50]byte // 5 keywords
	part     uint64      // i.e. part 1 of 5 parts
	parts    uint64
	geo      [2]byte // e.g. uk, us, etc
	data     [3840]byte
}

type Node struct {
	socketAddr      utils.SocketAddr
	knownHosts      []utils.SocketAddr  // find a way of locking this
	storage         map[uuid.UUID]Block // find away of locking this
	uptime          float64
	serverBehaviour func(*Node, utils.SocketAddr, []byte) []byte
	clientBehaviour func(*Node)
	routes          map[string]func(*Node, utils.SocketAddr, []byte) []byte

	//lock            sync.Mutex
}

// NewNode based on the local IP address of the computer, an OS allocated or specified port number and the desired
// memory usage. The max memory is specified in megabytes.
func NewNode(port uint16, maxMemory uint64, serverBehaviour func(*Node, utils.SocketAddr, []byte) []byte, clientBehaviour func(*Node)) (Node, error) {
	var node Node

	// Convert from mb to bytes
	maxMemoryInBytes := maxMemory * 1024 * 1024

	// check if max memory is more than some arbitrary min value (what is the minimum value that would be useful?)
	if maxMemory < 512 {
		return node, errors.New("allocated memory must be at least 512MB")
	} else if maxMemoryInBytes > memory.TotalMemory() {
		return node, errors.New("allocated memory must be less than the total system memory")
	}
	//else if maxMemoryInBytes > memory.FreeMemory() {
	//	return node, errors.New("allocated memory must be less than the free system memory")
	//}

	// Determine the preferred local ip of the machine
	localIpString := utils.GetOutboundIP()

	// Determine the capacity of the knownHosts list size based on user specified max memory
	maxKnownHosts := uint64((0.05 * float64(maxMemoryInBytes)) / utils.SocketAddressSize) // 5% of allocated memory is used for the known host list

	// Determine the upper limit of data block
	maxStorageBlocks := (maxMemoryInBytes - maxKnownHosts) / BlockSize // remaining memory is used for the data blocks

	node = Node{
		socketAddr:      utils.SocketAddr{Ip: localIpString, Port: port},
		knownHosts:      make([]utils.SocketAddr, 0, maxKnownHosts), // make a slice of known hosts of length and capacity maxKnownHosts
		storage:         make(map[uuid.UUID]Block, maxStorageBlocks),
		uptime:          0,
		serverBehaviour: serverBehaviour,
		clientBehaviour: clientBehaviour,
	}

	return node, nil
}

func (node *Node) StartNode() {
	// at the same time:
	// - call out for other nodes (multicast)
	// - generate thread-pool + start listening for connections and respond to them with the prescribed listening behaviour
	// - run client behaviour as prescribed
	go node.findPeer()
	go node.clientBehaviour(node)
	node.tcpListen()
}

func (node *Node) tcpListen() {
	// Create listener socket
	//node.lock.Lock()
	l, err := net.Listen("tcp", node.socketAddr.ToString())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes. (https://gobyexample.com/defer)
	defer l.Close()

	// Reassign the node's port to the actual port number of the TCP listener once it is created
	_, port, _ := net.SplitHostPort(l.Addr().String())
	portInt64, err := strconv.ParseUint(port, 10, 16)
	node.socketAddr.Port = uint16(portInt64)
	//node.lock.Unlock()

	fmt.Println("Node is listening at ", l.Addr())

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Received connection")
		// Handle connections in a new goroutine.
		go node.handleRequest(conn)
	}
}

func (node *Node) handleRequest(conn net.Conn) {
	defer conn.Close()

	// Make a buffer to hold incoming data.
	fmt.Println("handling the request")
	buf, err := ioutil.ReadAll(conn) // BUG: why is that blocking?
	if err == io.EOF {
		fmt.Println("Error reading:", err.Error())
		//os.Exit(1)
	}
	//var buf []byte
	//conn.Read(buf)
	fmt.Println("Received data: ", string(buf))

	// React appropriately to the incoming request
	// Check if the request matches any of the reserved routes roots (for internal working of the distributed system
	// else request handled by user defined server behaviour (which can have its own roots too)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr).IP
	remotePort := conn.RemoteAddr().(*net.TCPAddr).Port
	remoteSocketAddr := utils.SocketAddr{Ip: remoteAddr, Port: uint16(remotePort)}
	//fmt.Println("Remote address: ", remoteSocketAddr)
	res := node.routeHandler(buf, remoteSocketAddr)
	fmt.Println("Response: ", res)
	conn.Write(res)
}

// make routeHanlder always return something - always have confirmation
func (node *Node) routeHandler(packet []byte, remoteAddr utils.SocketAddr) []byte {
	// BUG: When it received that payload eiter the fmt.Fprint is messing with the payload or my parsePacket function
	// I can be fairly sure that bug is being cause by fmt.Fprint (in the introduceMyself function)
	uri, payload := utils.ParsePacket(packet)
	//fmt.Println("Received request to ", uri)
	switch uri {
	case AppCode:
		node.serverBehaviour(node, remoteAddr, payload)
		return []byte{success}
	//	TODO: Convert the uri human friendly strings to 1 byte code numbers - so they will always be the first byte in the packet way more efficient!!
	case pingCode:
		remoteHostAddress, _ := utils.FromJson(payload)
		//fmt.Println(remoteHostAddress.ToString())
		//node.lock.Lock()
		node.addNewKnownHost(remoteHostAddress)
		node.introduceMyself(remoteHostAddress)
		//node.lock.Unlock()
		return []byte{success}
	case helloCode: // TODO: stop ping and udp listening from here
		fmt.Println("cool now we know each other")
		remoteHostAddress, _ := utils.FromJson(payload) // BINGO! the bug comes from here - the payload is 1010 in length for some reason
		//node.lock.Lock()
		//fmt.Println(len(remoteHostAddress))
		//fmt.Println(remoteHostAddress)
		node.addNewKnownHost(remoteHostAddress)
		//node.lock.Unlock()
		return []byte{success}
	}
	return []byte{failure}
}

func (node *Node) introduceMyself(remoteHost utils.SocketAddr) {
	// Start a tcp client connection and send them my ip and port
	c, err := net.Dial("tcp", remoteHost.ToString())
	if err != nil {
		fmt.Println(err)
		return
	}
	nodeSocketAddress, _ := node.socketAddr.ToJson()
	c.Write(append([]byte{helloCode}, nodeSocketAddress...))
	c.Close()
}

func foundNodeHandler(src *net.UDPAddr, n int, b []byte, node *Node) {
	log.Println(n, "bytes read from", src)
	//packet := string(b[:n])
	packet := b[:n]
	//fmt.Println(packet)
	remoteAddr := utils.SocketAddr{
		Ip:   src.IP,
		Port: uint16(src.Port),
	}
	node.routeHandler(packet, remoteAddr)
}

// findPeer solves the cold start problem (many computers running but un-aware of each other)
func (node *Node) findPeer() {
	//If I get a multicast that isn't myself then add it to the known hosts and stop pinging and listening
	go discover.ListenForMulticasts(node, foundNodeHandler) // This should always be listening out for new nodes that might want to join the network
	discover.PingLAN(node)                                  // This should stop once it has found a peer
}

func (node *Node) addNewKnownHost(remoteHost utils.SocketAddr) (bool, error) {
	if len(node.knownHosts) < cap(node.knownHosts) {
		node.knownHosts = append(node.knownHosts, remoteHost)
		return true, nil
	}
	return false, errors.New("known hosts array is full")
}

func (node *Node) KnownHosts() []utils.SocketAddr {
	//node.lock.Lock()
	//defer node.lock.Unlock()
	return node.knownHosts
}

func ifErrorFail(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func (node *Node) createConnections(remoteHost utils.SocketAddr) net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteHost.ToString())
	ifErrorFail(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	ifErrorFail(err)
	return conn
}

func (node *Node) Request(remoteHost utils.SocketAddr, route []byte, payload []byte) []byte {
	conn := node.createConnections(remoteHost)
	defer conn.Close()

	response := make([]byte, 0)

	conn.Write(append(route, payload...))
	conn.Read(response)

	return response
}

func (node *Node) SocketAddr() utils.SocketAddr {
	return node.socketAddr
}

func (node *Node) registerRoute(route []byte, handler func([]byte, utils.SocketAddr)) {
	node.routes[route] = handler
}
