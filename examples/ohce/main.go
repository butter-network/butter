package main

import (
	"bufio"
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
	"os"
)

// This is a very simple example of a butter program: a reverse echo. A node sends a user specified message to each of
// it's known hosts, the hosts reply with the message reversed.

func send(remoteHost utils.SocketAddr, msg string) (string, error) {
	response, err := utils.Request(remoteHost, []byte("reverse-message/"), []byte(msg))
	if err != nil {
		return "", err
	}
	return string(response), nil
}

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
func serverBehaviour(node *node.Node, packet []byte) []byte {
	message := string(packet)
	reversedMsg := reverse(message)
	return []byte(reversedMsg)
}

// The clientBehavior for this application is to send a string to all the nodes known hosts for them to reverse it
func clientBehaviour(node *node.Node) {
	for {
		fmt.Print("Type message:")
		in := bufio.NewReader(os.Stdin)
		line, _ := in.ReadString('\n')

		knownHosts := node.KnownHosts()

		for i := 0; i < len(knownHosts); i++ {
			res, err := send(knownHosts[i], line)
			if err != nil {
				fmt.Println("unable to send message to", knownHosts[i])
			}
			fmt.Println(knownHosts[i].ToString(), "responded with: ", res)
		}
	}
}

func main() {
	// Create a new node by:
	// - Specifying a port or setting it to 0 to let the OS assign a port
	// - Defining an upper limit for the memory usage of the node on the system (recommended setting it to 2048mb)
	// - Specifying a clientBehaviour function to describe the interface for the user to interact with the decentralised application
	node, _ := node.NewNode(0, 2048, clientBehaviour)

	// Specifying a serverBehaviour function to be called when an app level packet is received
	node.RegisterRoute("reverse-message/", serverBehaviour)

	// Spawn your node into the butter network
	butter.Spawn(&node, false) // blocking
}
