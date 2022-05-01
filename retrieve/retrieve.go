// Package retrieve is Butter's inbuilt information retrieval solution. Implements a `naive' BFS, TTL BFS and RBFS
package retrieve

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/persist"
	"github.com/butter-network/butter/utils"
)

// When queried through the BFS mechanism
func retrieve(overlay node.Overlay, query []byte) []byte {
	persistOverlay := overlay.(*persist.Overlay)
	block, err := persistOverlay.Block(string(query))
	if err == nil {
		return append([]byte("found/"), block.Data()...)
	}

	addrs := make([]utils.SocketAddr, 0)
	// Add all my known hosts to the queue
	for host := range overlay.Node().KnownHosts() {
		addrs = append(addrs, host)
	}
	addrsJson, _ := json.Marshal(addrs)
	return append([]byte("try/"), addrsJson...)
}

// When queried through the RBFS mechanism
func rbfsretrieve(overlay node.Overlay, payload []byte) []byte {
	persistOverlay := overlay.(*persist.Overlay)
	// Separate the payload into the random node param and the query
	param, query, _ := utils.ParsePacket(payload)
	block, err := persistOverlay.Block(string(query))
	if err == nil {
		return append([]byte("found/"), block.Data()...)
	}

	// Make a random selection of known hosts and return those
	rand, _ := strconv.ParseFloat(string(param), 32)
	randomHosts := randomNodes(int(float32(rand)*float32(len(overlay.Node().KnownHosts()))), overlay.Node().KnownHosts())
	jsonHosts, _ := json.Marshal(randomHosts)
	return append([]byte("try/"), jsonHosts...)
}

func AppendRetrieveBehaviour(node *node.Node) {
	node.RegisterServerBehaviour("retrieve/", retrieve)
	node.RegisterServerBehaviour("random-bfs-retrieve/", rbfsretrieve)
}

// NaiveRetrieve High level entrypoint for searching for a specific piece of information on the network
// look if I have the information else look at the most likely known host to get to that information
// one query per piece of information (one-to-one) hence the query has to be unique i.e i.d.
func NaiveRetrieve(overlay persist.Overlay, query string) []byte {
	// do I have this information, if so return it
	// else BFS (pass the query on to all known hosts (partial view)
	block, err := overlay.Block(string(query))
	if err == nil {
		return block.Data()
	}
	return bfs(overlay, query)
}

func bfs(overlay persist.Overlay, query string) []byte {
	// Initialise an empty queue
	queue := make([]utils.SocketAddr, 10)
	// Add all my known hosts to the queue
	for host := range overlay.Node().KnownHosts() {
		queue = append(queue, host)
	}
	for len(queue) > 0 {
		// Pop the first element from the queue
		host := queue[0]
		queue = queue[1:]
		// Start a connection to the host, Ask host if he has data, receive response
		response, _ := utils.Request(host, []byte("retrieve/"), []byte(query))
		route, payload, err := utils.ParsePacket(response)
		if err != nil {
			fmt.Println("unable to parse packet")
		}
		// If the returned packet is success + the data then return it
		// else add the known hosts of the remote node to the end of the queue
		if string(route) == "found/" {
			return payload
		}
		// failed but gave us their known hosts to add to queue
		remoteKnownHosts, _ := utils.AddrSliceFromJson(payload)
		queue = append(queue, remoteKnownHosts...) // add the remote hosts to the end of the queue
	}
	return []byte("Information is not on the network")
}

func ttlBfs(overlay persist.Overlay, query string, ttl int) []byte {
	// Initialise an empty queue
	queue := make([]utils.SocketAddr, 0)
	// Add all my known hosts to the queue
	for host := range overlay.Node().KnownHosts() {
		queue = append(queue, host)
	}
	for len(queue) > 0 || ttl == 0 {
		// Pop the first element from the queue
		host := queue[0]
		queue = queue[1:]
		// Start a connection to the host, Ask host if he has data, receive response
		response, _ := utils.Request(host, []byte("retrieve/"), []byte(query))
		route, payload, err := utils.ParsePacket(response)
		if err != nil {
			fmt.Println("unable to parse packet")
		}
		// If the returned packet is success + the data then return it
		// else add the known hosts of the remote node to the end of the queue
		if string(route) == "found/" {
			return payload
		}
		// failed but gave us their known hosts to add to queue
		remoteKnownHosts, _ := utils.AddrSliceFromJson(payload)
		queue = append(queue, remoteKnownHosts...) // add the remote hosts to the end of the queue
		ttl--
	}
	return []byte("Information is not on the network")
}

func randomBfs(overlay persist.Overlay, query string, ttl int, prop float32) []byte {
	// Initialise an empty queue
	queue := make([]utils.SocketAddr, 0)
	// Add all my known hosts to the queue
	queue = append(queue, randomNodes(int(prop*float32(len(overlay.Node().KnownHosts()))), overlay.Node().KnownHosts())...)
	for len(queue) > 0 || ttl == 0 {
		// Pop the first element from the queue
		host := queue[0]
		queue = queue[1:]
		// Start a connection to the host, Ask host if he has data, receive response
		response, _ := utils.Request(host, []byte("retrieve/"), []byte(query))
		route, payload, err := utils.ParsePacket(response)
		if err != nil {
			fmt.Println("unable to parse packet")
		}
		// If the returned packet is success + the data then return it
		// else add the known hosts of the remote node to the end of the queue
		if string(route) == "found/" {
			return payload
		}
		// failed but gave us their known hosts to add to queue
		remoteKnownHosts, _ := utils.AddrSliceFromJson(payload)
		queue = append(queue, remoteKnownHosts...) // add the remote hosts to the end of the queue
		ttl--
	}
	return []byte("Information is not on the network")
}

// Return random set of nodes from KnownHosts
func randomNodes(n int, hosts map[utils.SocketAddr]node.HostQuality) []utils.SocketAddr {
	// select n random nodes from the list of hosts
	// copy to a new list
	// return the new list
	newHosts := make([]utils.SocketAddr, n)
	keys := make([]utils.SocketAddr, n)
	for key, _ := range hosts {
		keys = append(keys, key)
	}
	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(hosts))
		newHosts[i] = keys[randomIndex]
	}
	return newHosts
}
