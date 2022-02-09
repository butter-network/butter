package persist

import (
	"errors"
	"fmt"
	"github.com/a-shine/butter/node"
	uuid "github.com/nu7hatch/gouuid"
)

type Overlay struct {
	node    *node.Node
	storage map[uuid.UUID]Block
}

func (o *Overlay) Node() *node.Node {
	return o.node
}

// AddBlock to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (o *Overlay) AddBlock(keywords []string, data string) string {
	// TODO: add the logic to break down the data into blocks if it exceeds the block size
	id, _ := uuid.NewV4()
	o.storage[*id] = Block{
		keywords: naiveProcessKeywords(keywords),
		part:     1,
		parts:    1,
		data:     naiveProcessData(data),
	}
	return id.String()
}

// Block from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *Overlay) Block(id string) (Block, error) {
	parsedId, err := uuid.ParseHex(id)
	if err != nil {
		fmt.Println("Error parsing UUID:", err)
		return Block{}, err
	}
	fmt.Println("Parsed ID: ", parsedId)
	if val, ok := o.storage[*parsedId]; ok {
		return val, nil
	}
	return Block{}, errors.New("block not found")
}
