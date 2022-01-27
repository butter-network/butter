package store

import (
	"encoding/json"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
)

const (
	GetGroups                byte = 21
	GetGroupParticipantCount byte = 22
)

// Peer Content Group (PCG) logic

// Two levels of abstraction:
// - node level (network level)
// - group level (storage level) - think about groups and data/block interchangeably

// TODO: add groups to node struct - make group the storage property of node

type GroupPayload struct {
	uuid   [16]byte
	leader utils.SocketAddr
}

func (g *GroupPayload) getUuid() [16]byte {
	return g.uuid
}

type Group struct {
	// group ID which is inherent to the map of groups (no need to have it in the struct)
	uuid         [16]byte
	leader       utils.SocketAddr
	participants []utils.SocketAddr
	data         Block
}

type Block struct {
	//uuid     [16]byte // probably don't need this?
	//part     uint64   // i.e. part 1 of 5 parts
	//parts    uint64
	keywords [5][50]byte // 5 keywords
	geo      [2]byte     // e.g. uk, us, etc
	data     [3840]byte
}

func getGroupParticipantCount(node *node.Node, groupPayload GroupPayload) int {
	uri := node.GroupCode + GetGroupParticipantCount
	response := node.Request(groupPayload.leader, []byte{uri}, nil)
	return int(response[0])
}

// This function has a big overhead - cause basically you are looking at all the data of your known hosts
func findGroups(node *node.Node) []Group {
	knownHosts := node.KnownHosts()
	groupsSet := make(map[[16]byte]utils.SocketAddr) // map the group uuid to the leader
	groupOrderedSet := make([]byte, 0)

	// for each known host find the groups they are in
	// then query each group to see how many node are in the group
	// order list by number of nodes in group
	// based on node's available storage assign node to as many groups as possible (append to groups)

	for _, host := range knownHosts {
		hostGroups := node.Request(host, []byte{butter.GetGroups}, nil)
		// convert group json to grouppayload struct
		var groupsPayload []GroupPayload
		json.Unmarshal(hostGroups, &groupsPayload)
		// put groups into a set
		for _, group := range groupsPayload {
			groupsSet[group.getUuid()] = group.leader
		}
	}

	// query each group to see how many nodes are in the group
	for _, group := range groupsSet {
		groupParticipantNb := group.get(groupUUID)
		// TODO: order groupSet by number of nodes in group
	}
	return groupOrderedSet
}

// using as much redundancy as we can, and when we want to add something to the network we eat into that redundant data storage
// however there is an atomic group number (has to be at least 3 to allow for verification TMR redundancy model)
func join(node *node.Node) {
	//groups := findGroups(node)
	// for each group
	// - if group is found
	// - - send join request to members
	// - else create a group?
	orderedGroups := findGroups(node)

	for _, group := range orderedGroups {
		err := node.AddGroup(group)
		if err != nil {
			// I have no more capacity to join groups
			break
		}
	}
}

//func receiveMessageFromGroup(message []byte) {
//	switch message {
//	case "join":
//
//	}
//}
