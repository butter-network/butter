package main

import (
	"fmt"
	"github.com/a-shine/butter"
)

// This is a very simple example of a butter program: a reverse echo. A node sends a user specified message to each of
// it's known hosts, the hosts reply with the message reversed.

// function, which takes a string as
// argument and return the reverse of string.
func reverse(s string) string {
	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	// return the reversed string.
	return string(rns)
}

func serverBehaviour(node *butter.Node, incomingMsg []byte) []byte {
	incomingMsgString := string(incomingMsg)
	reversedMsg := reverse(incomingMsgString)
	return []byte(reversedMsg)
}

func clientBehaviour(node *butter.Node) {
	for {
		fmt.Println("Type message:")
		var msg string
		fmt.Scanln(&msg) // blocks until user input

		knownHosts := node.GetKnownHosts()

		for i := 0; i < len(knownHosts); i++ {
			res, err := butter.Send(knownHosts[i], msg)
			if err != nil {
				fmt.Println("unable to send message to", knownHosts[i])
			}
			fmt.Println(knownHosts[i], " responded with: ", res)
		}
	}
}

func main() {
	// Create a new node (define a port or set it to 0 to let the OS assign a port)
	// Define an upper limit of memory usage for the node on the system (recommended setting it to 2048mb (2GB)) or set to 0
	//to use all available memory
	node, err := butter.NewNode(0, 2048, serverBehaviour, clientBehaviour)
	if err != nil {
		fmt.Println(err)
		return
	}
	node.StartNode()
}
