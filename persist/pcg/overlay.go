package pcg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/butter-network/butter/node"
)

// Peer implements an overlay node, as described in the Butter node overlay interface.
type Peer struct {
	node           *node.Node
	maxStorage     uint64
	currentStorage uint64
	storage        map[[32]byte]*Group
}

// NewPCG constructs a Peer.
func NewPCG(node *node.Node, maxMemoryMb uint64) Peer {
	maxMemory := MbToBytes(maxMemoryMb)
	maxStorage := MaxStorage(maxMemory)
	return Peer{
		node:       node,
		maxStorage: maxStorage,
		storage:    make(map[[32]byte]*Group),
	}
}

// Node gets the underlying Node in the Peer object. Required by the Butter overlay interface.
func (p *Peer) Node() *node.Node {
	return p.node
}

// AvailableStorage gets the remaining storage available to a Peer. Required by the Butter overlay interface.
func (p *Peer) AvailableStorage() uint64 {
	return p.maxStorage - p.currentStorage
}

// Group gets a Peer's Group by its UUID. If the Peer is not a participant of the passed UUID's Group, returns nil with
// an error.
func (p *Peer) Group(id string) (*Group, error) {
	var hash [32]byte
	data, _ := hex.DecodeString(id)
	copy(hash[:], data)
	if group, ok := p.storage[hash]; ok {
		return group, nil
	}
	return nil, errors.New("block not found")
}

// Groups gets a Peer's Groups' UUIDs.
func (p *Peer) Groups() map[[32]byte]*Group {
	return p.storage
}

// CreateGroup creates a Group storing the passed data, and assigns the Peer as its first participant. The Group UUID is
// generated from the stored information. Returns the UUID of the new group as a string.
func (p *Peer) CreateGroup(data string) string {
	var formattedData [4096]byte
	copy(formattedData[:], data)
	hsha2 := sha256.Sum256(formattedData[:])
	p.storage[hsha2] = NewGroup(formattedData, p.node.SocketAddr())
	p.currentStorage += 4096
	return fmt.Sprintf("%x", hsha2)
}

// JoinGroup assigns the passed Group to a Peer's Group.
func (p *Peer) JoinGroup(g Group) {
	hsha2 := sha256.Sum256(g.Data[:])
	err := g.AddParticipant(p.node.SocketAddr())
	if err != nil {
		fmt.Println("Unable to join group:", err)
	}
	p.storage[hsha2] = &g
}

// String returns a string representation of a Peer.
func (p *Peer) String() string {
	str := ""
	for _, g := range p.Groups() {
		str = str + g.String()
	}
	return str
}
