package butter

import (
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/node"
)

func Spawn(node *node.Node, traverse bool) {
	go discover.Discover(node) // BUG: sometimes the node conn socket (created in the node.Listen() method) has not been created in time for discover to use the correct SocketAddr
	if traverse {
		//go traverse.Traverse(node)
	}
	go node.ClientBehaviour(node)
	node.Listen()
}
