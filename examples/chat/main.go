package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/utils"
)

// This example demonstrate how the library abstracts away much of the distributed networking so that teh app designer
// can focus on building functionality. In addition, it demonstrates the use of app level url routing.

// The serverBehaviour abstracts away all the distributed networking, so the app designer is only ever dealing with the
// app level packets
func serverBehaviour(node *butter.Node, remoteNodeAddr utils.SocketAddr, packet []byte) []byte {
	message := string(packet)
	fmt.Println(remoteNodeAddr.ToString()+" says: ", message)
	return []byte("/message-received")
}

// clientBehaviour creates an interface where a user can send a message to all his known hosts and get confirmation if
// they have received it.
func clientBehaviour(node *butter.Node) {
	for {
		fmt.Println("Type message:")
		var msg string
		fmt.Scanln(&msg)

		knownHosts := node.GetKnownHosts()

		for i := 0; i < len(knownHosts); i++ {
			response, err := butter.Send(knownHosts[i], msg)
			if err != nil {
				return
			}
			if string(response) == "/message-received" {
				fmt.Println(knownHosts[i].ToString() + " go the message successfully!")
			}
		}
	}
}

func main() {
	node, _ := butter.NewNode(0, 2048, serverBehaviour, clientBehaviour)
	node.StartNode()
}
