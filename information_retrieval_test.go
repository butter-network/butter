package butter

import (
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/retrieve"
	"github.com/a-shine/butter/store"
	"testing"
)

// TODO: Complete the test

const nodes = 10

var uuids = make([]string, 0)

func TestInformationRetrieval(t *testing.T) {
	// create many nodes and spawn into butter network
	for i := 0; i < nodes; i++ {
		node, _ := node.NewNode(0, 512, clientBehaviourTest, false)

		retrieve.AppendRetrieveBehaviour(&node)

		uuid := store.NaiveStore(&node)
		uuids = append(uuids, uuid)

		// Spawn your node into the butter network
		go Spawn(&node, false) // blocking
	}

	// add information to nodes (could be any random string)
	// store the uuid in a slice

	// then outside of loop create a new node
	node := node.NewNode(0, 512, clientBehaviour, false)
	retrieve.AppendRetrieveBehaviour(&node)
	Spawn(&node, false) // blocking
	// search for each piece of information by going over the uuids
	// time the amount of time taken to retrieve the information
}

func clientBehaviourTest(n *node.Node) {
	// do nothing
}

func clientBehaviour(n *node.Node) {
	// do nothing
	for i, uuid := range uuids {
		// start a timer
		retrieve.Retrieve(n, uuid)
		// stop timer
		if i == len(uuids)-1 {
			break
		}

	}
}
