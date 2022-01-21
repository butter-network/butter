package main

import (
	"fmt"
	"github.com/a-shine/butter"
)

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

func clientBehaviour(node *butter.Node) {
	for {
		fmt.Println("Type message:")
		var msg string
		fmt.Scanln(&msg) // blocks until user input

		knownHosts := node.GetKnownHosts()

		fmt.Println(knownHosts)

		for i := 0; i < len(knownHosts); i++ {
			//fmt.Println(len(knownHosts[i]))
			butter.Send(knownHosts[i], msg)
		}
	}
}

func serverBehaviour(node *butter.Node, incomingMsg string) string {
	return reverse(incomingMsg)
}

func main() {
	// Create a new node (define a port or set it to 0 to let the OS assign a port)
	// Define an upper limit of memory usage for the node on the system (recommended setting it to 2048mb (2GB)) or set to 0
	//to use all available memory
	node, err := butter.NewNode(0, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}
	node.StartNode(clientBehaviour, serverBehaviour)
}

// TODO: Fix the bug in the code, so that the chat works between several nodes
