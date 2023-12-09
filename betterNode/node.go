// Implementation of a 'better' node with an unspecified communication interface (like in libp2p). This node could
// communicate with varying communication protocols. This will be particularly important to improve the testbed so that
// pipes can be created for node communication instead of TCP connections. Implementation currently incomplete.

package betterNode

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/pbnjay/memory"
	"log"
	"net"
	"time"
)

const EOF byte = 26 // EOF code

type CommunicationInterface interface {
	Listen() (net.Listener, error)
	Connect(addr string) (net.Conn, error)
}

type Overlay interface {
	Node() *Node
	//	Add client and server behaviours specific to the overlay
	// add storage protocol specific to the overlay
}

type KnownHost struct {
	addr net.Addr
	// TODO: quality metadata for known host optimisation
}

type Node struct {
	id             uuid.UUID
	comsInterfaces []CommunicationInterface // Allow a node to communicate over multiple interfaces simultaneously
	started        time.Time
	overlayStack   []Overlay
}

func NewNode(commInterfaces []CommunicationInterface, allocatedMemoryMb uint64) (*Node, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Check that there is enough memory available to allocate to the node
	allocatedMemoryBytes := allocatedMemoryMb * 1024 * 1024
	availableMemoryBytes := memory.FreeMemory()
	if allocatedMemoryBytes > availableMemoryBytes {
		return nil, fmt.Errorf("allocated memory (%v bytes) is greater than available memory (%v bytes)", allocatedMemoryBytes, availableMemoryBytes)
	}

	// FIXME: Allocate 10% of allocated memory to storage of known host data
	//knownHostMem := uint64(0.1 * float64(allocatedMemoryBytes))

	node := &Node{
		id:             *u4,
		comsInterfaces: commInterfaces,
		overlayStack:   make([]Overlay, 0),
	}

	// persistent storage should be part of the overlay network, not the node

	return node, nil
}

func (n *Node) Start() {
	n.started = time.Now()

	// Start listening on all communication interfaces
	for _, comsInterface := range n.comsInterfaces {
		listener, err := comsInterface.Listen()
		if err != nil {
			fmt.Printf("Error starting listener: %v\n", err)
			continue
		}

		go n.handleConnections(listener)
	}

	// Start the node's default background behaviors
	// 1. Do peer discovery
}

func (n *Node) handleConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		log.Printf("Ready to accept and handle connections accepted on %v\n\n", listener.Addr().String())

		// Handle the connection in a separate goroutine
		go n.handleConnection(conn)
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	// Your code to handle communication over this connection goes here
	// Unpack the payload and determine the connection route/function to call
	// Call the appropriate function and send the response back to the client
}

func main() {
	// Create a new node with desired communication interfaces and known hosts
	// Initialize TCP communication interface
	tcpComm, err := NewTCPCommunication() // Specify your TCP listen address
	if err != nil {
		fmt.Printf("Error creating TCP communication: %v\n", err)
		return
	}

	// Initialize Pipe communication interface
	pipeComm, err := NewPipeCommunication() // Use a unique socket file path
	if err != nil {
		fmt.Printf("Error creating Pipe communication: %v\n", err)
		return
	}

	commInterfaces := []CommunicationInterface{
		tcpComm,
		pipeComm,
	}

	node, err := NewNode(commInterfaces, 2)
	if err != nil {
		fmt.Printf("Error creating node: %v\n", err)
		return
	}

	// Start the node
	node.Start()
}
