package butter

import "math/rand"

// graph struct
type graph struct {
	Nodes []node
}

// genRandomGraph generates a random graph with n nodes and m edges
func genRandomGraph(n, m int) graph {
	g := graph{}
	for i := 0; i < n; i++ {
		g.Nodes = append(g.Nodes, node{ID: i})
	}
	for i := 0; i < m; i++ {
		g.Nodes[rand.Intn(n)].Edges = append(g.Nodes[rand.Intn(n)].Edges, edge{ID: i})
	}
	return g
}

// -- Test retrieval performance --
// Randomly generate graphs
// Create them as Butter networks (where each vertex is a node and each edge is a link/tcp connection)
// Add a random number of files to each node
// Randomly select a file and a node to search for it
// Record time + number of nodes visited
