package store

import (
	"github.com/a-shine/butter/node"
)

// NaiveStore stores information on the network naively by simply placing it on the local node. It generate a UUIS for
// the information and creates an information block and return information uuid
func NaiveStore(node *node.Node, keywords []string, data string) string {
	uuid := node.AddBlock(keywords, data)
	return uuid
}
