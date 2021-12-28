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
		fmt.Scanln(&msg)

		//knownHosts := node.GetKnownHosts()
		//
		//for i := 0; i < knownHosts.len(); i++ {
		//	_, err := node.Send(knownHosts[i], msg)
		//	if err != nil {
		//		fmt.Println("Error sending message to ", knownHosts[i])
		//	}
		//}
	}
}

func serverBehaviour(node *butter.Node, incomingMsg string) string {
	return reverse(incomingMsg)
}

func main() {
	node := butter.NewNode(0)
	butter.StartNode(&node, clientBehaviour, serverBehaviour)
}
