package main

import (
	"fmt"
	"github.com/a-shine/butter"
)

func clientBehaviour(node *butter.Node) {
	for {
		var interactionType string
		fmt.Print("Would you like to add(1) or search(2) information on the network: ")
		fmt.Scanln(&interactionType)

		switch interactionType {
		case "1":
			var data string
			var keywords []string
			fmt.Println("What is the information you would like to store: ")
			fmt.Scanln(&data)
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
			butter.NaiveStore(node, keywords, data)
		case "2":
			var searchType string
			fmt.Println("Would you like to \n-retrieve(1) a specific piece of information or, \n-explore(2) information on the network:")
			fmt.Scanln(&searchType)
			switch searchType {
			case "1":
				var uuid string
				fmt.Println("What is the UUID of the piece of information you would like to retrieve: ")
				fmt.Scanln(&uuid)
				fmt.Println(butter.NaiveRetrieve(node, uuid))
			case "2":
			// TODO: implement search engine behaviour
			default:
				fmt.Println("Invalid choice")
			}
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func serverBehaviour(node *butter.Node, incomingMsg string) string {
	return ""
}

func main() {
	node := butter.NewNode(0)
	butter.StartNode(&node, clientBehaviour, serverBehaviour)
}
