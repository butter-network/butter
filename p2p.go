package butter

import (
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/traverse"
)

func Spawn(node *node.Node, traverseFlag bool) {
	go discover.Discover(node) // BUG: sometimes the node conn socket (created in the node.Listen() method) has not been created in time for discover to use the correct SocketAddr
	if traverseFlag {
		go traverse.Traverse(node)
	}
	node.Start()
}

// Embassadro has to be a traversed node (must port forward either manually or via UPNP) + extra embassadorial behaviour is added to node
func SpawnEmbassador(node *node.Node) {
	go discover.Discover(node) // BUG: sometimes the node conn socket (created in the node.Listen() method) has not been created in time for discover to use the correct SocketAddr
	go traverse.Traverse(node)
	//go traverse.AppendAmbassadorBehaviour(node)
	node.Start()
}

//func SimulationSpawn(node *node.Node, traverseFlag bool) {
//	go discover.SimulationDiscover(node)
//	node.SimulatedStart()
//}
