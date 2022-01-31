package butter

import (
	"github.com/a-shine/butter/node"
	"testing"
)

const nodes = 10

func TestInformationRetrieval(t *testing.T) {
	// create many nodes and spawn into butter network
	for i := 0; i < nodes; i++ {
		node, _ := node.NewNode(0, 512, clientBehaviour, false)

		// Spawn your node into the butter network
		Spawn(&node, false) // blocking
	}

	// add information to nodes (could be any random string)
	// store the uuid in a slice

	// then outside of loop create a new node
	// search for each piece of information by going over the uuids
	// time the amount of time taken to retrieve the information
}

func clientBehaviour(n *node.Node) {
	// do nothing
}
