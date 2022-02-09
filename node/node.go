package node

import (
	"errors"
	"fmt"
	"github.com/a-shine/butter/utils"
	"github.com/pbnjay/memory"
	"log"
	"net"
	"os"
	"strconv"
)

type Overlay interface {
	Node() *Node
}

type Node struct {
	listener        net.Listener
	knownHosts      []utils.SocketAddr // find a way of locking this
	uptime          float64
	ClientBehaviour func(Overlay)
	routes          map[string]func(Overlay, []byte) []byte
	simulated       bool
	ambassador      bool
}

// --- Getters ---

func (node *Node) Address() string {
	return node.listener.Addr().String()
}

func (node *Node) SocketAddr() utils.SocketAddr {
	socketAddr, _ := utils.AddrFromString(node.listener.Addr().String())
	return socketAddr
}

func (node *Node) KnownHosts() []utils.SocketAddr {
	return node.knownHosts
}

func (node *Node) KnownHostsStruct() utils.SocketAddrSlice {
	return node.knownHosts
}

func (node *Node) IsSimulated() bool {
	return node.simulated
}

// --- Setters ---

// UpdateIP of node to the given IP. This is important for updating an IP from local to public during NAT traversal.
func (node *Node) UpdateIP(ip string) {
	cachePort := node.SocketAddr().Port
	node.closeListener()
	node.listener, _ = net.Listen("tcp", ip+":"+strconv.Itoa(int(cachePort)))
}

// --- Adders ---

// RegisterRoute allows a node to register a behaviour for a route. This allows dapp designers to aff their own
// functionality and build on top of butter nodes. All routes have node and incoming payload as parameters and return a
// response payload.
func (node *Node) RegisterRoute(route string, handler func(Overlay, []byte) []byte) {
	node.routes[route] = handler
}

// AddKnownHost to increase node's partial view of the network. If already in known hosts, does nothing. If known
// hosts list is full, runs manageKnownHosts function (black box function that manages an optimal known host list).
func (node *Node) AddKnownHost(remoteHost utils.SocketAddr) {
	// TODO: check if the host is already known
	if len(node.knownHosts) < cap(node.knownHosts) {
		node.knownHosts = append(node.knownHosts, remoteHost)
	}
	node.manageKnownHosts()
}

// NewNode based on the local IP address of the computer, a port number, the desired memory usage and an application
// specific client behaviour. If the port is unspecified (i.e. 0), teh OS will allocate an available port. The max
// memory is specified in megabytes. If the memory is not specified (i.e. 0), the default is 2048 MB (2GB). A node has
// to contribute at least 512 MB of memory to the network (for it to be worthwhile) and use less memory than the total
// system memory.
func NewNode(port uint16, maxMemoryMb uint64, clientBehaviour func(Overlay), simulated bool) (Node, error) {
	var node Node

	// Sets the default memory to 2048 MB if not specified
	if maxMemoryMb == 0 {
		maxMemoryMb = 2048
	}

	// Convert user specified max memory in mb to bytes
	maxMemory := mbToBytes(maxMemoryMb)

	// check if max memory is more than some arbitrary min value (what is the minimum value that would be useful?)
	if maxMemory < mbToBytes(512) {
		return node, errors.New("allocated memory must be at least 512MB")
	} else if maxMemory > memory.TotalMemory() {
		return node, errors.New("allocated memory must be less than the total system memory")
	}

	// Determine the capacity of the knownHosts list size based on user specified max memory
	maxKnownHosts := uint64((0.05 * float64(maxMemory)) / float64(utils.SocketAddressSize)) // 5% of allocated memory is used for the known host list

	// Determine the upper limit of data block
	//maxStorageBlocks := (maxMemory - maxKnownHosts) / BlockSize // remaining memory is used for the data blocks

	//if simulated {
	//	// create a simulated listener?
	//	listener := mock_conn.NewConn()
	//} else {
	// Determine the preferred local ip of the machine
	localIpString := utils.GetOutboundIP() // TODO: Find a better way to do this

	var socketAddr utils.SocketAddr
	socketAddr.Ip = localIpString
	socketAddr.Port = port

	listener, err := net.Listen("tcp", socketAddr.ToString())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	//}

	node = Node{
		listener:        listener,
		knownHosts:      make([]utils.SocketAddr, 0, maxKnownHosts), // make a slice of known hosts of length and capacity maxKnownHosts
		uptime:          0,
		ClientBehaviour: clientBehaviour,
		routes:          make(map[string]func(Overlay, []byte) []byte),
		simulated:       simulated,
		ambassador:      false,
	}

	return node, nil
}

// Start node by listening out for incoming connections and starting the application specific client behaviour. A node
// behaves both as a server and a client simultaneously (that's how peer-to-peer systems work).
func (node *Node) Start(overlay Overlay) {
	go node.ClientBehaviour(overlay)
	node.listen(overlay)
}

func (node *Node) closeListener() {
	err := node.listener.Close()
	if err != nil {
		log.Println("Error closing listener:", err.Error())
		log.Println("Unable to shutdown gracefully")
	}
}

// Shutdown gracefully by closing the listener, telling the network the node is leaving and passing on data as required
func (node *Node) Shutdown() {
	// TODO: Gracefully shutdown method incomplete
	node.closeListener()
}

// listen to incoming connections from other nodes and handle them in serrate goroutines
func (node *Node) listen(overlay Overlay) {
	for {
		conn, err := node.listener.Accept()
		if err != nil {
			// Avoid fatal errors at all costs - we want to maximise node availability
			log.Println("Node is unable to accept incoming connections due to: ", err.Error())
			continue // forces next iteration of the loop skipping any code in between
		}

		// Pass connection to request handler in a new goroutine - allows a node to handle multiple connections at once
		go node.HandleRequest(conn, overlay)
	}
}

// HandleRequest by reading the connection buffer, processing the packet and writing the response to the connection
// buffer
func (node *Node) HandleRequest(conn net.Conn, overlay Overlay) {
	packet, err := utils.Read(&conn) // Read incoming buffer until EOF
	if err != nil {
		log.Println("Unable to read due to", err.Error())
	}

	// RouteHandler will handle invalid packet or route errors by returning an error uri. This allows us to always
	// handle requests without panicking.
	response := node.RouteHandler(packet, overlay)

	err = utils.Write(&conn, response)
	if err != nil {
		log.Println("Unable to write due to:", err.Error())
	}

	err = conn.Close()
	if err != nil {
		log.Println("Error closing connection:", err.Error())
	}
}

// RouteHandler for incoming packets that applies the correct response to the packet or returns an error
// ("invalid-packet/" if the node is unable to pass the packet or "invalid-route" if the node does not have a
// registered behaviour corresponding to the route)
func (node *Node) RouteHandler(packet []byte, overlay Overlay) []byte {
	route, payload, err := utils.ParsePacket(packet)
	if err != nil {
		return []byte("invalid-packet/")
	}

	// TODO: Don't think this works - need to test
	if response := node.routes[string(route)](overlay, payload); response != nil {
		return response
	}

	return []byte("invalid-route/")
}

func (node *Node) manageKnownHosts() {
	// TODO: Implement known hosts management
	// choose and maintain node host list
	// keep metadata about from previous node queries
	//learn about known hosts every time I deal with a request
	//make a known hosts list evaluator function
}

func mbToBytes(mb uint64) uint64 {
	return mb * 1024 * 1024
}
