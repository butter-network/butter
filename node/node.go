// Node is at the core of the butter framework. Every dapp built with the butter framework will be composed of many
// butter nodes. Each butter node will have behaviours that allow it to fulfill the functionality required from the
// dapp.

package node

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/a-shine/butter/utils"
	"github.com/pbnjay/memory"
)

// Overlay interface describes what an implemented Overlay struct should look like
type Overlay interface {
	Node() *Node
	AvailableStorage() uint64
}

type Node struct {
	listener         net.Listener
	quit             chan interface{}
	wg               sync.WaitGroup
	knownHosts       KnownHosts
	started          time.Time
	clientBehaviours []func(Overlay)
	serverBehaviours map[string]func(Overlay, []byte) []byte
	ambassador       bool
	storageMemoryCap uint64
}

// --- Getters ---

func (node *Node) Address() string {
	return node.listener.Addr().String()
}

func (node *Node) SocketAddr() utils.SocketAddr {
	socketAddr, _ := utils.AddrFromString(node.listener.Addr().String())
	return socketAddr
}

func (node *Node) KnownHosts() map[utils.SocketAddr]HostQuality {
	return node.knownHosts.Addrs()
}

func (node *Node) KnownHostsStruct() KnownHosts {
	return node.knownHosts
}

func (node *Node) StorageMemoryCap() uint64 {
	return node.storageMemoryCap
}

// --- Setters ---

// UpdateIP of node to the given IP. This is important for updating an IP from local to public during NAT traversal.
func (node *Node) UpdateIP(ip string) {
	cachePort := node.SocketAddr().Port
	node.closeListener()
	node.listener, _ = net.Listen("tcp", ip+":"+strconv.Itoa(int(cachePort)))
}

// --- Adders ---

// A node has a collection of behaviours that determine its functionality. Client behaviours determine how a user might
// be expected to interact with a node while severe behaviours determine how a node might respond to certain requests.

// RegisterClientBehaviour allows you to define several ways to interface with a node. The ability to append
// user-defined functions allows you to add functionality to the node specific to your dapp.
func (node *Node) RegisterClientBehaviour(handler func(Overlay)) {
	node.clientBehaviours = append(node.clientBehaviours, handler)
}

// RegisterServerBehaviour allows a node to register a behaviour for a route to increase its ability to respond too
// requests. All routes (specific server handler functions) take the overlay network and incoming payload as parameters,
// process the request and return a response payload.
func (node *Node) RegisterServerBehaviour(route string, handler func(Overlay, []byte) []byte) {
	node.serverBehaviours[route] = handler
}

// AddKnownHost to increase node's partial view of the network. If already in known hosts, does nothing. If known
// hosts list is full, determines intelligently if the host should be added and which should be removed. Be careful, a
// node is not always added despite the function call.
func (node *Node) AddKnownHost(remoteHost utils.SocketAddr) {
	node.knownHosts.Add(remoteHost)
}

// RemoveKnownHost from node's list. If the host is not in the list, does nothing.
func (node *Node) RemoveKnownHost(remoteHost utils.SocketAddr) {
	node.knownHosts.Remove(remoteHost)
}

// --- Constructor ---

// NewNode based on the local IP address of the computer, a port number, the desired memory usage. If the port is
// unspecified (i.e. 0), the OS will allocate an available port. The max memory is specified in megabytes. If the memory
// is not specified (i.e. 0), the default is 512 MB (0.5GB). A node has to contribute at least 512 MB of memory to the
// network (for it to be worthwhile) and use less memory than the total system memory.
func NewNode(port uint16, maxMemoryMb uint64) (*Node, error) {
	var node Node

	// Sets the default memory to 512 MB if not specified
	if maxMemoryMb == 0 {
		maxMemoryMb = 512
	}

	// Convert user specified max memory in mb to bytes
	maxMemory := mbToBytes(maxMemoryMb)

	// check if max memory is more than some arbitrary min value (what is the minimum value that would be useful?)
	if maxMemory < mbToBytes(512) {
		return &node, errors.New("allocated memory must be at least 512MB")
	} else if maxMemory > memory.TotalMemory() {
		return &node, errors.New("allocated memory must be less than the total system memory")
	}

	// Determine the capacity of the KnownHosts list size based on user specified max memory
	knownHostsMemory := uint64(0.0001 * float64(maxMemory)) // 0.001% of allocated memory is used for the known host list ~10,000 known hosts for 512MB
	knownHostsCap := int(knownHostsMemory) / utils.SocketAddressSize

	// Determine the upper limit of storage in bytes (so that the overlay network has an idea of how much memory it can use)
	maxStorage := maxMemory - knownHostsMemory // remaining memory is used for the storage

	ip, _ := utils.GetIp()

	var socketAddr utils.SocketAddr
	socketAddr.Ip = ip
	socketAddr.Port = port

	listener, err := net.Listen("tcp", socketAddr.ToString())
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	node = Node{
		listener:         listener,
		quit:             make(chan interface{}),
		knownHosts:       KnownHosts{cap: uint(knownHostsCap), Hosts: make(map[utils.SocketAddr]HostQuality)}, // make a slice of known hosts of length and capacity knownHostsCap
		started:          time.Time{},
		clientBehaviours: make([]func(Overlay), 0),
		serverBehaviours: make(map[string]func(Overlay, []byte) []byte),
		ambassador:       false,
		storageMemoryCap: maxStorage,
	}

	node.wg.Add(1)

	return &node, nil
}

// Start node by listening out for incoming connections and starting the application specific client behaviours. A node
// behaves both as a server and a client simultaneously (that's how peer-to-peer systems work!).
func (node *Node) Start(overlay Overlay) {
	node.RegisterClientBehaviour(updateKnownHosts) // Periodically remove dead hosts from known hosts list
	AppendHostQualityServerBehaviour(node)         // Append host quality server behaviour to the node
	node.started = time.Now()
	for i := range node.clientBehaviours {
		go node.clientBehaviours[i](overlay)
	}
	node.listen(overlay)
}

// TODO: Wait for all connections to close before returning
func (node *Node) closeListener() {
	close(node.quit)
	err := node.listener.Close()
	if err != nil {
		log.Println("Error closing listener:", err.Error())
		log.Println("Unable to shutdown gracefully")
	}
	node.wg.Wait()
}

// Shutdown gracefully by closing the listener, telling the network the node is leaving and passing on data as required
func (node *Node) Shutdown() {
	// TODO: Gracefully shutdown method incomplete
	// TODO: add shoutdown methoth to overlay network so that it can be called from here - this would allow overlay network designers to implement their own graceful shutdown behaviour
	node.closeListener()
}

// listen to incoming connections from other nodes and handle them in serrate goroutines
func (node *Node) listen(overlay Overlay) {
	defer node.wg.Done()

	for {
		conn, err := node.listener.Accept()
		if err != nil {
			select {
			case <-node.quit:
				return
			default:
				// Avoid fatal errors at all costs - we want to maximise node availability
				log.Println("Node is unable to accept incoming connections due to: ", err.Error())
				//continue // forces next iteration of the loop skipping any code in between
			}
		} else {
			node.wg.Add(1)
			go func() {
				// Pass connection to request handler in a new goroutine - allows a node to handle multiple connections at once
				node.HandleRequest(conn, overlay)
				node.wg.Done()
			}()
		}
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
	if response := node.serverBehaviours[string(route)](overlay, payload); response != nil {
		return response
	}

	return []byte("invalid-route/")
}

// uptime of the node i.e. time since it started listening and hence contributing to the butter network
func (node *Node) uptime() time.Duration {
	return time.Since(node.started)
}

// updateKnownHosts by periodically (every 2 min) querying each known host for its HostQuality and storing the updated
// information in the KnownHosts cache.
func updateKnownHosts(overlay Overlay) {
	for {
		overlay.Node().knownHosts.update() // updates host metadata + removes dead hosts
		time.Sleep(time.Second * 120)
	}
}
