package pcg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/utils"
	// "github.com/a-shine/butter/node"
	// "github.com/a-shine/butter/utils"
)

// Storage overlay route APIs.
const (
	inGroupUri = "in-group?/"
	canJoinUri = "can-join?/"
)

var alreadyFinding bool // flag to check if the leader is already finding new participants for the group

// Store adds the passed data to the network using the PCG protocol. Returns the UUID for the data on the network.
func Store(overlay *Peer, data string) string {
	// Creates Group for passed Peer storing passed data
	uuid := overlay.CreateGroup(data)
	return uuid
}

// AppendGroupStoreBehaviour registers the behaviours that allow the node to work with the PCG overlay
func AppendGroupStoreBehaviour(node *node.Node) {
	node.RegisterServerBehaviour(inGroupUri, inGroup)
	node.RegisterServerBehaviour(canJoinUri, canJoin)
	node.RegisterClientBehaviour(heartbeat)
}

// inGroup is a server behaviour - a querying node can ask a given node if it is in a group
func inGroup(overlayInterface node.Overlay, groupId []byte) []byte {
	pcg := overlayInterface.(*Peer)
	_, err := pcg.Group(string(groupId))
	if err != nil {
		return []byte("Group not found")
	}
	return []byte("Group member")
}

// canJoin is a server behaviour - a querying node can ask a given node if it has the memory capacity to join a group
func canJoin(overlayInterface node.Overlay, payload []byte) []byte {
	peer := overlayInterface.(*Peer)
	if peer.currentStorage < peer.maxStorage {
		var groupDigest Group
		err := json.Unmarshal(payload, &groupDigest)
		if err != nil {
			fmt.Println("error marshaling group")
		}
		peer.JoinGroup(groupDigest)
		return []byte("accepted")
	}
	return []byte("can't join group")
}

// heartbeat is a client behaviour (always running in the background of the node) that schedules a Node to manage its
// Groups' participants
func heartbeat(overlayInterface node.Overlay) {
	pcgn := overlayInterface.(*Peer)
	for {
		manageParticipants(pcgn)
		time.Sleep(time.Second * 2)
	}
}

// amILeader returns whether a given Peer is the leader of a given Group
func (p *Peer) amILeader(g *Group) bool {
	socketAddr := p.Node().SocketAddr()
	socketAddrStr := socketAddr.ToString()
	if !GroupContains(g.Participants, socketAddr) {
		return false
	}
	for _, h := range g.Participants {
		if h.ToString() > socketAddrStr {
			return false
		}
	}
	return true
}

// manageParticipants triggers a given Peer to manage its Groups' participants. For every Group of Peer and Participant
// of Group, check if Participant is still in Group, and update Participants list accordingly. If the Group has fewer
// than 3 Participants, find new replacement Participants.
func manageParticipants(peer *Peer) {
	for id, group := range peer.Groups() {
		for _, participant := range group.Participants {
			response, err := utils.Request(participant, []byte(inGroupUri), id[:])
			if err != nil || string(response) != "Group not found" {
				err := group.RemoveParticipant(participant)
				if err != nil {
					fmt.Println("Error removing participant:", err)
				}
			}
		}
		if peer.amILeader(group) && ((len(group.Participants)) < 3) && !alreadyFinding {
			go findParticipants(peer, group)
		}
	}
}

// GroupContains checks if a host is withing a group's participants
func GroupContains(g []utils.SocketAddr, h utils.SocketAddr) bool {
	for _, a := range g {
		if a.ToString() == h.ToString() {
			return true
		}
	}
	return false
}

// findParticipants finds Peers from a given Peer's known hosts to assign to a Group. This is done by the leader of the
// group.
func findParticipants(pcg *Peer, group *Group) {
	// Checks that they don't already belong to the Group.
	// Checks that they have enough storage for the Group's data.
	// Breaks out of while loop once Group has the correct number of Participants.
	alreadyFinding = true
	for {
		for host, _ := range pcg.Node().KnownHosts() {
			if GroupContains(group.Participants, host) {
				continue
			}
			output, err := json.Marshal(group)
			if err != nil {
				break
			}
			response, err := utils.Request(host, []byte(canJoinUri), output)
			if err != nil || string(response) == "no storage available" {
				fmt.Println(err)
			}
			if string(response) == "accepted" {
				err := group.AddParticipant(host)
				if err != nil {
					fmt.Println(err)
				}
				if len(group.Participants) == 3 {
					break
				}
			}
		}
		time.Sleep(time.Second * 1)
		if len(group.Participants) == 3 {
			break
		}
	}
	alreadyFinding = false
}
