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

const BlockSize = 4096

// A Block has a size of 4096 bytes with a uuid of size 16 bytes, 5 keywords of max size 50 bytes, 2 part numbers of
// size 8 bytes each, a geotag of size 2 bytes and a data of size 4096 - 16 - 5*50 - 2*8 - 2 = 3840 bytes. A Block is
// uniquely identified by combining its uuid and part number e.g. <UUID>/<PartNumber>.
type Block struct {
	keywords [5][50]byte // 5 keywords
	part     uint64      // i.e. part 1 of 5 parts
	parts    uint64
	geo      [2]byte // e.g. uk, us, etc
	data     [3840]byte
}

func (b *Block) Data() []byte {
	return b.data[:]
}

type Node struct {
	listener        net.Listener
	knownHosts      []utils.SocketAddr  // find a way of locking this
	storage         map[uuid.UUID]Block // find away of locking this
	uptime          float64
	ClientBehaviour func(*Node)
	routes          map[string]func(*Node, []byte) []byte
	simulated       bool
}

// NewNode based on the local IP address of the computer, an OS allocated or specified port number and the desired
// memory usage. The max memory is specified in megabytes.
func NewNode(port uint16, maxMemory uint64, clientBehaviour func(*Node), simulated bool) (Node, error) {
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

	// Determine the capacity of the knownHosts list size based on user specified max memory
	maxKnownHosts := uint64((0.05 * float64(maxMemoryInBytes)) / float64(utils.SocketAddressSize)) // 5% of allocated memory is used for the known host list

	// Determine the upper limit of data block
	maxStorageBlocks := (maxMemoryInBytes - maxKnownHosts) / BlockSize // remaining memory is used for the data blocks

	//if simulated {
	//	// create a simulated listener?
	//	listener := mock_conn.NewConn()
	//} else {
	// Determine the preferred local ip of the machine
	localIpString := utils.GetOutboundIP()

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
	route, payload := utils.ParsePacket(packet)
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

// just store as much of the data as possible - cut off the rest
func naiveProcessData(data string) [3840]byte {
	var formattedData [3840]byte
	for i, _ := range formattedData {
		formattedData[i] = data[i]
	}
	return formattedData
}

func processKeywords(keywords []string) [5][50]byte {
	var formattedKeywords [5][50]byte
	for i, _ := range formattedKeywords {
		var word [50]byte
		for j, _ := range word {
			formattedKeywords[i][j] = keywords[i][j]
		}
	}
	return formattedKeywords
}

func (node *Node) IsSimulated() bool {
	return node.simulated
}

func (node *Node) KnownHostsStruct() utils.SocketAddrSlice {
	return node.knownHosts
}

func (node *Node) Address() string {
	return node.listener.Addr().String()
}
