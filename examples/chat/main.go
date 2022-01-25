package main

import (
	"fmt"
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/utils"
	"net"
)

// This example demonstrate how the library abstracts away much of the distributed networking so that teh app designer
// can focus on building functionality. In addition, it demonstrates the use of app level url routing.

func send(remoteHost utils.SocketAddr, payload []byte) {
	c, err := net.Dial("tcp", remoteHost.ToString())
	if err != nil {
		fmt.Println(err)
		//return nil, errors.New("could not connect to remote host")
	}

	// Append the payload to the appCode to create the packet to send
	//var eof byte = 26
	packet := append([]byte{butter.AppCode}, payload...) // appCode is for app level requests
	response := make([]byte, 0)
	//packet = append(packet, io.EOF)
	c.Write(packet)
	c.Read(response)
	c.Close()
	fmt.Println("res: " + string(response))
	//response, err := ioutil.ReadAll(c)
	//fmt.Fprint(c, string(packet))

	if err != nil {
		fmt.Println(err)
		//return nil, errors.New("could not read response from remote host")
	}
	//if response == "/message-received" {
	//	c.Close()
	//}

	//fmt.Fprint(c, message)
	//response, _ := ioutil.ReadAll(c)
	//fmt.Printf(string(response))
	//return response, nil // TODO: fix this design This is blocking now
}

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

		knownHosts := node.KnownHosts()

		for i := 0; i < len(knownHosts); i++ {
			send(knownHosts[i], []byte(msg))
		}
	}
}

func main() {
	node, _ := butter.NewNode(0, 2048, serverBehaviour, clientBehaviour)
	node.StartNode()
}
