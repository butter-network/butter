package main

import (
	"github.com/a-shine/butter"
	"github.com/a-shine/butter/discover"
	"github.com/a-shine/butter/retrieve"
	"github.com/a-shine/butter/store"
	"github.com/a-shine/butter/utils"
)

func publish(node *Node, article string) {
	store.Store(node, []byte(article))
}

func read(node *Node, uuid [16]byte) string {
	article, err := retrieve.Item(node, uuid)
	if err != nil {
		return "unable to retrieve article"
	}
	return string(article)
}

// Interesting to think about but ignore this functionality for the moment
//func listAllArticles(node *Node) []string {
//	articles, err := retrieve.SearchDomain(node, []byte("/blog"))
//	if err != nil {
//		return []string{"unable to list articles"}
//	}
//	return articles
//}

//func listMyArticles(node *Node) []string {
//	articles, err := retrieve.SearchDomain(node, []byte("/blog/alex"))
//	if err != nil {
//		return []string{"unable to list articles"}
//	}
//	return articles
//}

func clientBehaviour(node *butter.Node) {
	// publish an article
	// read an article
}

func serverBehaviour(node *butter.Node, remoteHost utils.SocketAddr, appPacket []byte) []byte {
	return []byte("hello world!")
}

func main() {
	node, _ := butter.NewNode(0, 2048, serverBehaviour, clientBehaviour) // non-blocking
	discover.Discover(*node)
	//traverse.Traverse(*node) // this is not required but if you want the node to traverse nat this is required (traverse update teh node socket address to be a public IP)
	node.Start()

}
