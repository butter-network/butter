package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
)

// This is a very simple example of a butter program: a reverse echo. A node sends a user specified message to each of
// it's known hosts, the hosts reply with the message reversed.

func send(remoteHost utils.SocketAddr, msg string) (string, error) {
	response, err := utils.Request(remoteHost, []byte("message/"), []byte(msg))
	if err != nil {
		return "", err
	}
	return string(response), nil
}

// The serverBehavior for this application is to reverse the packets it receives and send them back to the sender
func serverBehaviour(node *node.Node, packet []byte) []byte {
	message := string(packet)
	fmt.Println("Message received: ", message)
	return []byte("success")
}

// The clientBehavior for this application is to send a string to all the nodes known hosts for them to reverse it
func clientBehaviour(node *node.Node) {
	for {
		fmt.Println("Type message:")
		var msg string
		fmt.Scanln(&msg) // Blocks until user input

		knownHosts := node.KnownHosts()

		for i := 0; i < len(knownHosts); i++ {
			fmt.Println("Sending message to:", knownHosts[i])
			res, err := send(knownHosts[i], msg)
			fmt.Println("Received:", res)
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
	// - Specifying a clientBehaviour function to describe the interface for the user to interact with the decentralised application
	node, _ := node.NewNode(0, 2048, clientBehaviour)

	// Specifying a serverBehaviour function to be called when an app level packet is received
	node.RegisterRoute("message/", serverBehaviour)

	// Spawn your node into the butter network
	butter.Spawn(&node, false) // blocking
}
