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
