package store

import "github.com/a-shine/butter"

// Peer Content Group (PCG) logic

// Two levels of abstraction:
// - node level (network level)
// - group level (storage level) - think about groups and data interchangably

// TODO: add groups to node struct - make group the storage property of node

func findGroups(node *butter.Node) [][]byte {
	//knownHosts := node.getKnownHosts()
	groups := make([][]byte, 0)

	// for each known host find the groups they are in
	// then query each group to see how many node are in the group
	// order list by number of nodes in group
	// based on node's available storage assign node to as many groups as possible (append to groups)

	return groups
}

// using as much redundancy as we can, and when we want to add something to the network we eat into that redundant data storage
// however there is an atomic group number (has to be at least 3 to allow for verification TMR redundancy model)

func join(node *butter.Node) {
	//groups := findGroups(node)
	// for each group
	// - if group is found
	// - - send join request to members
	// - else create a group?
}

//func receiveMessageFromGroup(message []byte) {
//	switch message {
//	case "join":
//
//	}
//}
