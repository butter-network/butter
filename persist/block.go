package persist

import (
	"unsafe"
)

const BlockSize = uint64(unsafe.Sizeof(Block{}))

// A Block is the atomic unit of storage for butter. Extra meta-data is attached to the data field (i.e. keywords,
// part-count and geotag) to improve storage and IR performance. A Block is uniquely identified by combining its uuid
//	and part number e.g. <UUID>/<PartNumber>.
type Block struct {
	keywords [5][50]byte // 5 keywords
	part     uint64      // i.e. part 1 of 5 parts
	parts    uint64
	geo      [2]byte // e.g. uk, us, etc
	data     [3840]byte
}

// Data field getter for a Block
func (b *Block) Data() []byte {
	return b.data[:]
}
