// Package persist is one of the Butter persist overlay implementations. Other are available on the butter-network
// GitHub repository.
package persist

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/butter-network/butter/utils"
	"strconv"

	"github.com/butter-network/butter/node"
)

// Overlay that complies with the Butter persist interface
type Overlay struct {
	node       *node.Node
	storageCap uint64
	storage    map[Id]Block
}

type Id struct {
	hash [32]byte
	part int
}

func stringToId(s string) (Id, error) {
	var id Id
	hash, part, err := utils.ParsePacket([]byte(s))
	if err != nil {
		return id, errors.New("invalid id")
	}
	var decodedHash [32]byte
	data, _ := hex.DecodeString(string(hash))
	copy(decodedHash[:], data)
	id.hash = decodedHash
	id.part, _ = strconv.Atoi(string(part))
	return id, nil
}

func (o *Overlay) Node() *node.Node {
	return o.node
}

func (o *Overlay) AvailableStorage() uint64 {
	usedStorage := len(o.storage) * BlockSize
	return o.storageCap - uint64(usedStorage)
}

// AddBlock to the node's storage. A UUID is generated for every bit of information added to the network (no update
// functionality yet!). Returns the UUID of the new block as a string.
func (o *Overlay) addBlock(id Id, keywords [5][50]byte, data [4096]byte, parts int) {
	o.storage[id] = Block{
		keywords: keywords,
		parts:    parts,
		data:     data,
	}
}

// Block from the node's storage by its UUID. If the block is not found, an empty block with an error is returned.
func (o *Overlay) Block(idString string) (Block, error) {
	id, _ := stringToId(idString)
	if val, ok := o.storage[id]; ok {
		return val, nil
	}
	return Block{}, errors.New("block not found")
}

func NewOverlay(node *node.Node) Overlay {
	storageCap := node.StorageMemoryCap() / uint64(BlockSize) // infer the storage cap from the nodes allocated memory for storage
	return Overlay{
		node:       node,
		storageCap: storageCap,
		storage:    make(map[Id]Block),
	}
}

func (o *Overlay) AddInformation(keywords []string, data []byte) string {
	hash := sha256.Sum256(data)
	// Append the part nb to the hash - so once we have found one block, we can determine the other hashes we need to
	// find - allows us to parallelize because we don't need to wait for a block to find the next one finding all the
	// blocks we need to find
	chunks := chunking(data)
	keywordsFormatted := naiveProcessKeywords(keywords)
	for i, chunk := range chunks {
		part := i + 1
		id := Id{hash, part}
		o.addBlock(id, keywordsFormatted, chunk, len(chunks)) // distribute the blocks across nodes - don't naiveStore the entirety of a piece of data on one node but spread it out (like what Adam was saying similar to RAID)
	}
	return fmt.Sprintf("%x", hash) // encode the hash as hex string
}

// chunking information into smaller pieces that can be maintained by Blocks
func chunking(data []byte) [][4096]byte {
	var chunks [][4096]byte
	for i := 0; i < len(data); i += 4096 {
		chunk := [4096]byte{}
		end := i + 4096

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(data) {
			end = len(data)
		}

		copy(chunk[:], data[i:end])

		chunks = append(chunks, chunk)
	}

	return chunks
}
