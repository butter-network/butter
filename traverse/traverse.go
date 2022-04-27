package traverse

import (
	"fmt"
	"github.com/butter-network/butter/node"
	"gitlab.com/NebulousLabs/go-upnp"
	"log"
)

// update the listener to listen on the public IP address of the network + make sure the user has port forwarded the same port as the listener
// add node to ambassdors list

func Traverse(node *node.Node) {
	withUPNP(node)
}

func withUPNP(node *node.Node) {
	portForwardWithUPNP(node)
	addNodeToKnownAmbassadors(node)
}

func portForwardWithUPNP(node *node.Node) {
	// connect to router
	d, err := upnp.Discover()
	if err != nil {
		log.Fatal(err)
	}

	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Your external IP is:", ip)

	// forward a port
	err = d.Forward(node.SocketAddr().Port, "upnp test")
	if err != nil {
		log.Fatal(err)
	}

	// update node to use external IP
	node.UpdateIP(ip)

	//// un-forward a port
	//err = d.Clear(9001)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// record router's location
	//loc := d.Location()
	//
	//// connect to router directly
	//d, err = upnp.Load(loc)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func addNodeToKnownAmbassadors(node *node.Node) {
	// ask local known hosts if they know an ambassador node
	// send ambassador your public IP and forwarded port
}
