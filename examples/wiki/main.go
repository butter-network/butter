// Example of a butter dapp (decentralised application) where data is persistent: wiki. The basic functionality of the
// wiki is to be able to add an entry and read an entry.
package main

import (
	"bufio"
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/persist"
	"github.com/a-shine/butter/persist/retrieve"
	"github.com/a-shine/butter/persist/store"
	"os"
)

func addArticle(overlay *persist.Overlay) {
	var keywords []string
	fmt.Println("What is the information you would like to store: ")
	in := bufio.NewReader(os.Stdin)
	data, _ := in.ReadString('\n') // Read string up to newline
	fmt.Println("What keywords would you like to associate with this information: ")
	for {
		fmt.Print("Add a keyword (or press enter to quit): ")
		var keyword string
		fmt.Scanln(&keyword)
		if keyword == "" {
			break
		}
		keywords = append(keywords, keyword)
	}
	articleUuid := store.NaiveStore(overlay, keywords, data)
	fmt.Println("Your article has been stored with UUID: ", articleUuid)
}

func readArticle(overlay *persist.Overlay) {
	var searchType string
	fmt.Println("Would you like to \n-retrieve(1) a specific piece of information or, \n-explore(2) information on the network:")
	fmt.Scanln(&searchType)
	switch searchType {
	case "1":
		var uuid string
		fmt.Println("What is the UUID of the piece of information you would like to retrieve: ")
		fmt.Scanln(&uuid)
		fmt.Println(string(retrieve.NaiveRetrieve(overlay, uuid)))
	case "2":
	// TODO: implement search engine behaviour
	default:
		fmt.Println("Invalid choice")
	}
}

func clientBehaviour(appInterface interface{}) {
	overlay := appInterface.(*persist.Overlay) // uses the included persist package to describe the overlay network
	for {
		var interactionType string
		fmt.Print("Would you like to add(1) or search(2) information on the network: ")
		fmt.Scanln(&interactionType)

		switch interactionType {
		case "1":
			addArticle(overlay)
		case "2":
			readArticle(overlay)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func main() {
	// Create a new node by: specifying a port (or setting it to 0 to let the OS assign one), defining an upper limit on
	// memory usage (recommended setting it to 2048mb) and specifying a clientBehaviour function that describes the
	// user-interface to interact with the decentralised application
	butterNode, _ := node.NewNode(0, 2048, clientBehaviour, false)

	fmt.Println("Node is listening at", butterNode.Address())

	// No need to specify retrieval or storage server behaviours as they are handled by the provided butter storage and
	//retrieve packages

	overlay := persist.NewOverlay(&butterNode) // Creates a new overlay network

	// Spawn your node into the butter network
	butter.Spawn(&butterNode, overlay, false) // Blocking
}
