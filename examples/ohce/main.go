// Example of a butter dapp (decentralised application) where data is not persistent: reverse echo. A node sends a user
// specified message to each of it's known hosts, the hosts reply with the same message reversed.
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

// Takes as input a string and returns the string in reverse.
func reverse(s string) string {
	rns := []rune(s) // Convert string to rune array
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		// Swap the letters of the string
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

// The serverBehavior for this application is to reverse the packet it receives and return it back to the sender as a
// response
func serverBehaviour(appInterface node.Overlay, payload []byte) []byte {
	_ = appInterface.(*OverlayNode)
	message := string(payload)
	reversedMsg := reverse(message)
	return []byte(reversedMsg)
}

// send a message to a specified host via the application specified reverse-message/ route
func send(remoteHost utils.SocketAddr, msg string) (string, error) {
	response, err := utils.Request(remoteHost, []byte("reverse-message/"), []byte(msg)) // Uses the utils package (recommended)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

// The clientBehavior for this application is to send a string to all the node's known hosts and ask them to reverse it
// and return it back
func clientBehaviour(appInterface node.Overlay) {
	peer := appInterface.(*OverlayNode)
	// Create an input loop
	for {
		fmt.Print("Type message:")
		in := bufio.NewReader(os.Stdin)
		line, _ := in.ReadString('\n') // Read string up to newline

		knownHosts := peer.node.KnownHosts() // Get the node's known hosts

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
	butterNode, _ := node.NewNode(0, 2048, false)

	fmt.Println("Node is listening at", butterNode.Address())

	// Specifying app level server behaviours - you can specify as many as you like as long as they are not reserved by
	// other butter packages
	butterNode.RegisterServerBehaviour("reverse-message/", serverBehaviour) // The client behaviour interacts with this route
	butterNode.RegisterClientBehaviour(clientBehaviour)

	// Spawn your node into the butter network
	butter.Spawn(&OverlayNode{node: &butterNode}, false) // Blocking
}
