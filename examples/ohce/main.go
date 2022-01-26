package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/utils"
)

func send(remoteHost utils.SocketAddr, msg string) (string, error) {
	response, err := utils.Request(remoteHost, []byte("/reverse-message"), []byte(msg))
	if err != nil {
		return "", err
	}
	return string(response), nil

}

// This is a very simple example of a butter program: a reverse echo. A node sends a user specified message to each of
// it's known hosts, the hosts reply with the message reversed.

// Takes as input a string and returns the string in reverse
func reverse(s string) string {
	rns := []rune(s) // Convert string to rune array
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		// Swap the letters of the string
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

// The serverBehavior for this application is to reverse the packets it receives and send them back to the sender
func serverBehaviour(node *butter.Node, packet []byte) []byte {
	message := string(packet)
	reversedMsg := reverse(message)
	return []byte(reversedMsg)
}

// The clientBehavior for this application is to send a string to all the nodes known hosts for them to reverse it
func clientBehaviour(node *butter.Node) {
	for {
		fmt.Println("Type message:")
		var msg string
		fmt.Scanln(&msg) // Blocks until user input

		knownHosts := node.KnownHosts()

		for i := 0; i < len(knownHosts); i++ {
			res, err := send(knownHosts[i], msg)
			if err != nil {
				fmt.Println("unable to send message to", knownHosts[i])
			}
			fmt.Println(knownHosts[i], " responded with: ", res)
		}
	}
}

func main() {
	// Create a new node by:
	// - Specifying a port or setting it to 0 to let the OS assign a port
	// - Defining an upper limit for the memory usage of the node on the system (recommended setting it to 2048mb)
	// - Specifying a serverBehaviour function to be called when an app level packet is received
	// - Specifying a clientBehaviour function to describe the interface for the user to interact with the decentralised application
	node, _ := butter.NewNode(0, 2048, clientBehaviour) // non-blocking
	discover.Discover(&node)
	//traverse.Traverse(&node) // this is not required but if you want the node to traverse nat this is required (traverse update teh node socket address to be a public IP
	node.RegisterRoute("/reverse-message", serverBehaviour)
	node.StartNode()

}
