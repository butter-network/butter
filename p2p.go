package butter

import (
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/node"
)

func Spawn(node *node.Node, traverse bool) {
	go discover.Discover(node)
	if traverse {
		//go traverse.Traverse(node)
	}
	go node.ClientBehaviour(node)
	node.Listen()
}
