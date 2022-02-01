package butter

import (
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/retrieve"
	"github.com/a-shine/butter/traverse"
)

// Spawn node into the network (the node serves as an entry-point to the butter network). You can also do this manually
// to have more control over the specific protocols used in your dapp. This function presents a simple abstraction with
// the included default butter protocols.
func Spawn(node *node.Node, traverseFlag bool) {
	go discover.Discover(node)
	if traverseFlag {
		go traverse.Traverse(node)
	}
	retrieve.AppendRetrieveBehaviour(node)
	node.Start()
}

// SpawnAmbassador node which is a special community node with added ambassadorial behaviours that help it bridge
// connections across subnetworks. To be an ambassador a node inherently needs to be available publicly (must port
// forward either manually or via UPNP and have a public IP address). The added ambassadorial behaviours allows the node
// to share the public addresses of other traversed (i.e. public) nodes between each other.
func SpawnAmbassador(node *node.Node) {
	go discover.Discover(node)
	go traverse.Traverse(node)
	//go traverse.AppendAmbassadorBehaviour(node) // the node keeps track of ambassador so if someone needs an ambassador they can find them dynamically (improvement on bootstrapping)
	node.Start()
}
