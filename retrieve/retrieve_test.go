package retrieve

import (
	"fmt"
	"github.com/butter-network/butter"
	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/store"
	"testing"
	"time"
)

// TODO: Complete the test

const nodes = 10

var uuids = make([]string, 0)

func TestInformationRetrieval(t *testing.T) {
	// create many nodes and spawn into butter network
	for i := 0; i < nodes-1; i++ {
		n, _ := node.NewNode(0, 512)

		fmt.Println("Node created -", n.Address())

		dummyKeywords := []string{"dummy", "dummy", "dummy", "dummy", "dummy"}
		uuid := store.NaiveStore(n, dummyKeywords, "dummy")
		uuids = append(uuids, uuid)

		// Spawn your node into the butter network
		go butter.SpawnDefaultOverlay(n, false) // blocking
		fmt.Println(n.KnownHosts())
	}

	// add information to nodes (could be any random string)
	// store the uuid in a slice

	time.Sleep(time.Second * 5)

	// then outside of loop create a new node
	n, _ := node.NewNode(0, 512)
	fmt.Println("Node created -", n.Address())
	butter.SpawnDefaultOverlay(n, false) // blocking
	// search for each piece of information by going over the uuids
	// time the amount of time taken to retrieve the information
}

func dummyClientBehaviour(n *node.Node) {
	// do nothing - just using them for their listener functionality do noy need to interact with them
}

func clientBehaviour(n *node.Node) {
	fmt.Println(n.KnownHosts())
	for _, uuid := range uuids {
		// start a timer
		data := NaiveRetrieve(n, uuid)
		fmt.Println(string(data))
		// stop timer
	}
}
