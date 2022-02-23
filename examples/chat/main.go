// Example of a butter dapp (decentralised application) where data is not persistent: chat. A node sends a message to
// all its known hosts. The recipient nodes print the message to console for user to read.
package main

import (
	"bufio"
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
	"os"
)

type OverlayNode struct {
	node *node.Node
	// You can any other fields you might need to create an overlay network...
}

func (n *OverlayNode) Node() *node.Node {
	return n.node
}

func (n *OverlayNode) AvailableStorage() uint64 {
	return 0
}

// The serverBehavior for this application is to print the received message to console for the user to read and return
// a confirmation receipt
func serverBehaviour(_ node.Overlay, payload []byte) []byte {
	message := string(payload)
	fmt.Println("Received:", message)
	return []byte("received/")
}

// send a message to a specified host via the application specified message/ route
func send(remoteHost utils.SocketAddr, msg string) (string, error) {
	response, err := utils.Request(remoteHost, []byte("message/"), []byte(msg)) // Uses the utils package (recommended)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

// The clientBehavior for this application is to send a string to all the node's known hosts and see if they have
// received it successfully
func clientBehaviour(overlayInterface node.Overlay) {
	overlay := overlayInterface.(*OverlayNode)
	// Create an input loop
	for {
		fmt.Print("Type message:")
		in := bufio.NewReader(os.Stdin)
		line, _ := in.ReadString('\n') // Read string up to newline

		knownHosts := overlay.node.KnownHosts() // Get the node's known hosts

		for host := range knownHosts { // For each known host
			res, err := send(host, line) // Ask them to reverse the input message
			if err != nil {
				// If there is an error, log the error BUT DO NOT FAIL - in decentralised application we avoid fatal
				// errors at all costs as we want to maximise node availability
				fmt.Println("Unable to send message to", host)
			}
			fmt.Println(host.ToString(), "responded with:", res)
		}
	}
}

func main() {
	// Create a new node by: specifying a port (or setting it to 0 to let the OS assign one), defining an upper limit on
	// memory usage (recommended setting it to 2048mb) and specifying a clientBehaviour function that describes the
	// user-interface to interact with the decentralised application
	butterNode, _ := node.NewNode(0, 2048)

	fmt.Println("Node is listening at", butterNode.Address())

	// Specifying app level server behaviours - you can specify as many as you like as long as they are not reserved by
	// other butter packages
	butterNode.RegisterServerBehaviour("message/", serverBehaviour) // The client behaviour interacts with this route
	butterNode.RegisterClientBehaviour(clientBehaviour)

	// Spawn your node into the butter network
	butter.Spawn(&OverlayNode{node: butterNode}, false) // Blocking
}
