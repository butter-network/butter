package retrieve

import (
	"fmt"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/persist"
	"github.com/a-shine/butter/utils"
)

func retrieve(overlay *persist.Overlay, query []byte) []byte {
	block, err := overlay.Block(string(query))
	if err == nil {
		return append([]byte("found/"), block.Data()...)
	}

	hostsStruct := overlay.Node().KnownHostsStruct()
	knownHostsJson, _ := hostsStruct.ToJson()
	return append([]byte("try/"), knownHostsJson...)
}

func found(node *node.Node, query []byte) []byte {
	return query
}

func try(node *node.Node, query []byte) []byte {
	return query
}

func AppendRetrieveBehaviour(node *node.Node) {
	node.RegisterRoute("retrieve/", retrieve)
	//node.RegisterRoute("found/", found)
	//node.RegisterRoute("try/", try)
}

// NaiveRetrieve High level entrypoint for searching for a specific piece of information on the network
// look if I have the information else look at the most likely known host to get to that information
// one query per piece of information (one-to-one) hence the query has to be unique i.e i.d.
func NaiveRetrieve(overlay *persist.Overlay, query string) []byte {
	// do I have this information, if so return it
	// else BFS (pass the query on to all known hosts (partial view)
	block, err := overlay.Block(string(query))
	if err == nil {
		return block.Data()
	}
	return bfs(overlay, query)
}

func bfs(overlay *persist.Overlay, query string) []byte {
	// Initialise an empty queue
	queue := make([]utils.SocketAddr, 0)
	// Add all my known hosts to the queue
	for _, host := range overlay.Node().KnownHosts() {
		queue = append(queue, host)
	}
	for len(queue) > 0 {
		// Pop the first element from the queue
		host := queue[0]
		queue = queue[1:]
		// Start a connection to the host, Ask host if he has data, receive resposnse
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
