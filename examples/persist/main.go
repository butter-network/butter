// Demo CLI decentralised application where data persists beyond the existence of a particular node instance. Run the
// demo in several terminal instances (creates several nodes), add and retrieve information by the interfacing with the
// different nodes. If you kill a node, notice that you are still able to retrieve information that was added by that
// node.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/butter-network/butter"
	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/store/pcg"
)

// clear the terminal so that the interface is clean and easier to read
func clear() {
	fmt.Print("\033[H\033[2J")
}

// add information to the network by calling the pcg.Store function
func add(overlay *pcg.Peer) {
	fmt.Println("Input information:")
	in := bufio.NewReader(os.Stdin)
	data, _ := in.ReadString('\n') // Read string up to newline
	uuid := pcg.Store(overlay, strings.TrimSpace(data))
	clear()
	fmt.Println("UUID:", uuid)
	fmt.Println("Data:", strings.TrimSpace(data))
	fmt.Println("Enter to continue...")
	_, err := in.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	clear()
}

// retrieve information from the network by calling the pcg.Retrieve function
func retrieve(overlay *pcg.Peer) {
	fmt.Println("Information UUID:")
	in := bufio.NewReader(os.Stdin)
	uuid, _ := in.ReadString('\n') // Read string up to newline
	data, _ := pcg.NaiveRetrieve(overlay, strings.TrimSpace(uuid))
	clear()
	fmt.Println(string(data))
	fmt.Println("Enter to continue...")
	_, err := in.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	clear()
}

// printAll the groups that the node is currently a member of i.e. what data is stored on this node
func printAll(peer *pcg.Peer) {
	fmt.Println(peer.String())
	fmt.Println("Enter to continue...")
	in := bufio.NewReader(os.Stdin)
	_, err := in.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	clear()
}

// interact client behaviour that allows a user to interface with the network
func interact(overlayInterface node.Overlay) {
	peer := overlayInterface.(*pcg.Peer)
	fmt.Println("Node created!")
	fmt.Println("Socket address:", peer.Node().SocketAddr())
	time.Sleep(2 * time.Second)
	clear()
	for {
		// prompt to pcgStore or pcgRetrieve information
		var interactionType string
		fmt.Print("Add(1), Retrieve(2) information or List my groups(3)?")
		fmt.Scanln(&interactionType)
		clear()
		switch interactionType {
		case "1":
			add(peer)
		case "2":
			retrieve(peer)
		case "3":
			printAll(peer)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func main() {
	// Creates a new Butter node with an OS allocated port and 512MB of memory
	butterNode, _ := node.NewNode(0, 512)

	// Register the client behaviour
	butterNode.RegisterClientBehaviour(interact)

	// Initialise the PCG overlay
	overlay := pcg.NewPCG(butterNode, 512)
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	// Spawn the node into the network with the appended demo client behaviour and pcg overlay protocols
	butter.Spawn(&overlay, false, false)
}
