package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/a-shine/butter/node"
)

type Overlay struct {
	node    *node.Node
	storage map[[32]byte]Block
}

func (o *Overlay) Node() *node.Node {
	return o.node
}

// AddBlock to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (o *Overlay) AddBlock(keywords []string, data string) string {
	// TODO: add the logic to break down the data into blocks if it exceeds the block size
	hsha2 := sha256.Sum256([]byte(data))
	o.storage[hsha2] = Block{
		keywords: naiveProcessKeywords(keywords),
		part:     1,
		parts:    1,
		data:     naiveProcessData(data),
	}
	return fmt.Sprintf("%x", hsha2)
}

// Block from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *Overlay) Block(id string) (Block, error) {
	var hash [32]byte
	data, _ := hex.DecodeString(id)
	copy(hash[:], data)
	if val, ok := o.storage[hash]; ok {
		return val, nil
	}
	return Block{}, errors.New("block not found")
}

func NewOverlay(node *node.Node) Overlay {
	return Overlay{
		node:    node,
		storage: make(map[[32]byte]Block),
	}
}
