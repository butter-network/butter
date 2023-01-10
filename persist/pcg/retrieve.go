// Pulled and modified from Butter's original implementation of Information Retrieval
// https://github.com/a-shine/butter/blob/main/retrieve/retrieve.go (commit 20ffb299fb196bfe0386ee8ab02987b0fc5e0119)

package pcg

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/utils"
)

const foundEndpoint = "found/"
const tryEndpoint = "try/"

// retrieve behaviour for a PCG node. When queried, it will either return the information if it is part of the group
// responsible for hosting it, else it will return its known hosts so that the querying node can continue querying the
// network.
func retrieve(overlay node.Overlay, query []byte) []byte {
	persistOverlay := overlay.(*Peer)

	// Check if node has data (i.e. is part of group)
	group, err := persistOverlay.Group(string(query))
	if err == nil {
		return append([]byte(foundEndpoint), group.Data[:]...) // queryHit
	}

	// Otherwise, if data is not found, return byte array of known hosts to allow for further search
	addrs := make([]utils.SocketAddr, 0)

	// Add all my known hosts to the queue
	for host := range overlay.Node().KnownHosts() {
		addrs = append(addrs, host)
	}
	addrsJson, _ := json.Marshal(addrs)

	return append([]byte(tryEndpoint), addrsJson...)
}

// AppendRetrieveBehaviour to the Butter node (much like registering an HTTP route in a tradition backend web framework)
func AppendRetrieveBehaviour(node *node.Node) {
	node.RegisterServerBehaviour("pcgRetrieve/", retrieve)
}

// NaiveRetrieve entrypoint to search for a specific piece of information on the network by UUID (information hash)
// using a 'naive' BFS approach
func NaiveRetrieve(overlay *Peer, query string) ([]byte, error) {
	// One query per piece of information (one-to-one) hence the query has to be a unique id

	// Do I have this information, if so return it
	// else BFS (pass the query on to all known hosts)
	block, err := overlay.Group(query)
	if err == nil {
		return block.Data[:], nil
	}
	return bfs(overlay, query)
}

// bfs across the network until information is found. This is not particularly well suited to production and won't scale
// well. However, for testing it provides a deterministic means of checking if information exists on the network.
func bfs(overlay *Peer, query string) ([]byte, error) {
	// Initialise an empty queue
	queue := make([]utils.SocketAddr, 0)

	// Keep track of already visited nodes - map of checked nodes
	checked := make(map[utils.SocketAddr]bool)

	// Add all my known hosts to the queue
	for host := range overlay.Node().KnownHosts() {
		queue = append(queue, host)
	}

	for {
		// If the queue is empty we have explored all connected nodes and have not found the information
		if len(queue) <= 0 {
			break
		}

		host := queue[0] // Pop the first element from the queue
		queue = queue[1:]
		checked[host] = true // Mark as checked

		// Start a connection to the host - ask host if he has data
		response, err := utils.Request(host, []byte("pcgRetrieve/"), []byte(query))
		if err != nil {
			fmt.Println("error in request")
			fmt.Println(err)
			continue
		}

		// Parse response
		route, payload, err := utils.ParsePacket(response)
		if err != nil {
			fmt.Println("unable to parse packet")
			fmt.Println(err)
			continue
		}

		// If the returned packet is success then we have found the information
		// else add the known hosts of the remote node to the end of the queue
		if string(route) == foundEndpoint {
			return payload, nil
		} else if string(route) == tryEndpoint {
			// Failed but gave us their known hosts to add to queue
			remoteKnownHosts, _ := utils.AddrSliceFromJson(payload)

			// Iterate through new know hosts and only add to queue if not already checked
			for _, x := range remoteKnownHosts {
				if !checked[x] {
					queue = append(queue, x)
				}
			}
		} else {
			log.Println("Unexpected response")
		}
	}

	return []byte(""), errors.New("failed to retrieve information")
}
