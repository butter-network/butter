package pcg

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/butter-network/butter/utils"
)

const GroupStructSize = unsafe.Sizeof(Group{})

const DataReplicationCount = 3

// ParticipantCount is an alias for the DataReplicationCount i.e. you can think of data replication and group
// participants interchangeably
const ParticipantCount = DataReplicationCount

// A Group is a set of participants and the data that they are responsible for maintaining.
type Group struct {
	Participants []utils.SocketAddr
	Data         [4096]byte
}

// NewGroup constructs a Group.
func NewGroup(data [4096]byte, participant utils.SocketAddr) *Group {
	return &Group{
		Participants: []utils.SocketAddr{participant},
		Data:         data,
	}
}

// SetParticipants assigns nodes to be Group participants.
func (g *Group) SetParticipants(participants []utils.SocketAddr) {
	g.Participants = participants
}

// AddParticipant assigns a node to be a Group participant.
func (g *Group) AddParticipant(host utils.SocketAddr) error {
	if len(g.Participants) >= ParticipantCount {
		return errors.New("group is full")
	}
	g.SetParticipants(append(g.Participants, host))
	return nil
}

// RemoveParticipant removes a participant node from a Group
func (g *Group) RemoveParticipant(host utils.SocketAddr) error {
	for i, participant := range g.Participants {
		if participant.ToString() == host.ToString() {
			g.Participants = append(g.Participants[:i], g.Participants[i+1:]...)
			break
		}
		if i == len(g.Participants) {
			return errors.New("host not in group")
		}
	}
	return nil
}

// ToJson returns a JSON representation of the group.
func (g *Group) ToJson() []byte {
	groupJson, _ := json.Marshal(g)
	return groupJson
}

// String returns a string representation of the group.
func (g *Group) String() string {
	return fmt.Sprintf("Data: %s\nGroup Members: %v\nUUID: %x\n\n", g.Data[:], g.Participants, sha256.Sum256(g.Data[:]))
}
