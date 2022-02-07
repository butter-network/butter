package persist

import (
	"errors"
	"fmt"
	"github.com/a-shine/butter/node"
	uuid "github.com/nu7hatch/gouuid"
)

type Overlay struct {
	node    *node.Node
	routes  map[string]func(*Overlay, []byte) []byte // overlay level routes
	storage map[uuid.UUID]Block                      // Introduce the notion of persistent storage expressed in the Overlay network
}

func (o *Overlay) Node() *node.Node {
	return o.node
}

func NewOverlay(node *node.Node) Overlay {
	return Overlay{
		node:    node,
		storage: make(map[uuid.UUID]Block),
	}
}

func (o *Overlay) RegisterRoute(route string, handler func(*Overlay, []byte) []byte) {
	o.routes[route] = handler
}

// Determine the upper limit of data block
//maxStorageBlocks := (maxMemory - maxKnownHosts) / persist.BlockSize // remaining memory is used for the data blocks - maybe do this based on node?

// Block from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (node *Overlay) Block(id string) (Block, error) {
	parsedId, err := uuid.ParseHex(id)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return Block{}, err
	}
	fmt.Println("Parsed ID: ", parsedId)
	if val, ok := node.storage[*parsedId]; ok {
		return val, nil
	}
	return Block{}, errors.New("block not found")
}

// AddBlock to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (node *Overlay) AddBlock(keywords []string, data string) string {
	// TODO: add the logic to break down the data into blocks if it exceeds the block size
	id, _ := uuid.NewV4()
	node.storage[*id] = Block{
		keywords: naiveProcessKeywords(keywords),
		part:     1,
		parts:    1,
		data:     naiveProcessData(data),
	}
	return id.String()
}
