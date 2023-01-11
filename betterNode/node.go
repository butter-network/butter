// Implementation of a 'better' node with an unspecified communication interface (like in libp2p). This node could
// communicate with varying communication protocols. This will be particularly important to improve the testbed so that
// pipes can be created for node communication instead of TCP connections. Implementation currently incomplete.

package betterNode

import (
	"github.com/butter-network/butter/node"
	uuid "github.com/nu7hatch/gouuid"
	"time"
)

const EOF byte = 26 // EOF code

type Listener interface {
	Accept() (Conn, error)
	Close() error
	Addr() Addr
}

type Conn interface {
	read()
	write()
}

type CommunicationInterface interface {
	Listen()                                                        // blocking - continuously open channel listening for incoming connections
	Request(CommunicationInterface, []byte, []byte) ([]byte, error) // called every time a new request needs to be made
}

type Overlay interface {
	Node() *Node
	AvailableStorage() uint64
}

type Node struct {
	id               uuid.UUID
	commInterface    CommunicationInterface // could make this a slice and allow a node to communicate over many interfaces simultaneously
	knownHosts       node.KnownHosts
	started          time.Time
	clientBehaviours []func(Overlay)                         // can only access the Request() method
	serverBehaviours map[string]func(Overlay, []byte) []byte // can only access the Request() method
	ambassador       bool
	storageMemoryCap uint64
}

func NewNode(commInterface CommunicationInterface) (Node, error) {
	var node Node

	u4, err := uuid.NewV4()
	if err != nil {
		return node, err
	}

	node.id = *u4
	node.commInterface = commInterface
	node.clientBehaviours = make([]func(Overlay), 0)
	node.serverBehaviours = make(map[string]func(Overlay, []byte) []byte)

	return node, nil
}

func (n *Node) Start() {
	n.started = time.Now()

	n.commInterface.Listen()
}
