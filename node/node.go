package node

import (
	"errors"
	"fmt"
	"github.com/a-shine/butter/utils"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/pbnjay/memory"
	"net"
	"os"
)

type Node struct {
	listener        net.Listener
	knownHosts      []utils.SocketAddr  // find a way of locking this
	storage         map[uuid.UUID]Block // find away of locking this
	uptime          float64
	ClientBehaviour func(*Node)
	routes          map[string]func(*Node, []byte) []byte
	simulated       bool
	ambassador      bool
}

// NewNode based on the local IP address of the computer, a port number, the desired memory usage and an application
// specific client behaviour. If the port is unspecified (i.e. 0), teh OS will allocate an available port. The max
// memory is specified in megabytes. If the memory is not specified (i.e. 0), the default is 2048 MB (2GB). A node has
// to contribute at least 512 MB of memory to the network (for it to be worthwhile) and use less memory than the total
// system memory.
func NewNode(port uint16, maxMemoryMB uint64, clientBehaviour func(*Node), simulated bool) (Node, error) {
	var node Node

	// Sets the default memory to 2048 MB if not specified
	if maxMemoryMB == 0 {
		maxMemoryMB = 2048
	}

	// Convert from mb to bytes
	maxMemory := mbToBytes(maxMemoryMB)

	// check if max memory is more than some arbitrary min value (what is the minimum value that would be useful?)
	if maxMemory < mbToBytes(512) {
		return node, errors.New("allocated memory must be at least 512MB")
	} else if maxMemory > memory.TotalMemory() {
		return node, errors.New("allocated memory must be less than the total system memory")
	}

	// Determine the capacity of the knownHosts list size based on user specified max memory
	maxKnownHosts := uint64((0.05 * float64(maxMemory)) / float64(utils.SocketAddressSize)) // 5% of allocated memory is used for the known host list

	// Determine the upper limit of data block
	maxStorageBlocks := (maxMemory - maxKnownHosts) / BlockSize // remaining memory is used for the data blocks

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
		storage:         make(map[uuid.UUID]Block, maxStorageBlocks),
		uptime:          0,
		ClientBehaviour: clientBehaviour,
		routes:          make(map[string]func(*Node, []byte) []byte),
		simulated:       simulated,
		ambassador:      false,
	}

	return node, nil
}

func (node *Node) Start() {
	go node.listen()
	node.ClientBehaviour(node)
}

func (node *Node) listen() {
	for {
		// Listen for an incoming connection.
		conn, err := node.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Received connection")
		defer conn.Close()
		// Handle connections in a new goroutine.
		go node.newHandleRequest(conn)
	}
}

func (node *Node) newHandleRequest(conn net.Conn) {
	packet, err := utils.Read(conn)
	if err != nil {
		return
	}
	response := node.NewRouteHandler(packet) // handle invalid route error but do not panic - just ignore
	fmt.Println("Response: ", response)
	utils.Write(conn, response)
}

func (node *Node) NewRouteHandler(packet []byte) []byte { //TODO return invalid route error
	fmt.Println(string(packet))
	route, payload, _ := utils.ParsePacket(packet)
	fmt.Println("Received request to ", string(route))
	fmt.Println("Payload: ", string(payload))
	response := node.routes[string(route)](node, payload)
	return response
}

func (node *Node) AddNewKnownHost(remoteHost utils.SocketAddr) (bool, error) {
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

func (node *Node) SocketAddr() utils.SocketAddr {
	socketAddr, _ := utils.AddrFromString(node.listener.Addr().String())
	return socketAddr
}

func (node *Node) RegisterRoute(route string, handler func(*Node, []byte) []byte) {
	node.routes[route] = handler
}

// choose and maintain node host list
// keep metadata about from previous node queries

func (node *Node) shutdown() {
	node.listener.Close()
	// pass data on to someone else
}

func (node *Node) UpdateIP(ip string) {
	node.listener.Close()
	keepPort := node.SocketAddr().Port
	node.listener, _ = net.Listen("tcp", ip+":"+string(keepPort))
}

func manageKnownHosts(node *Node) {
	//learn about known hosts every time I deal with a request
	//make a known hosts list evaluator function
}

func (node *Node) GetBlock(id string) (Block, error) {
	parsedId, _ := uuid.Parse([]byte(id))
	if val, ok := node.storage[*parsedId]; ok {
		return val, nil
	}
	return Block{}, errors.New("block not found")
}

func (node *Node) AddBlock(keywords []string, data string) string {
	// potentially add the logic to break down the data into it's component parts
	id, _ := uuid.NewV4()
	node.storage[*id] = Block{
		keywords: processKeywords(keywords),
		part:     1,
		parts:    1,
		data:     naiveProcessData(data),
	}
	return id.String()
}

//// just store as much of the data as possible - cut off the rest
//func naiveProcessData(data string) [3840]byte {
//	var formattedData [3840]byte
//	for i, _ := range formattedData {
//		formattedData[i] = data[i]
//	}
//	return formattedData
//}
//
//func processKeywords(keywords []string) [5][50]byte {
//	var formattedKeywords [5][50]byte
//	for i, _ := range formattedKeywords {
//		var word [50]byte
//		for j, _ := range word {
//			formattedKeywords[i][j] = keywords[i][j]
//		}
//	}
//	return formattedKeywords
//}

func (node *Node) IsSimulated() bool {
	return node.simulated
}

func (node *Node) KnownHostsStruct() utils.SocketAddrSlice {
	return node.knownHosts
}

func (node *Node) Address() string {
	return node.listener.Addr().String()
}
