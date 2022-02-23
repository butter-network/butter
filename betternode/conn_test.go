package betternode

import (
	"testing"
)

func TestPipeCommunication(t *testing.T) {
	node1, _ := NewNode()
	node2, _ := NewNode()

	go node1.Start()
	go node2.Start()

	response, _ := node1.commInterface.Request(node2.commInterface, []byte("ping/"), []byte(""))
	if string(response) != "pong" {
		t.Errorf("Expected pong, got %s", string(response))
	}
}

func TestTCPCommunication(t *testing.T) {
	node1, _ := NewNode()
	node2, _ := NewNode()

	go node1.Start()
	go node2.Start()

	response, _ := node1.commInterface.Request(node2.commInterface, []byte("ping/"), []byte(""))
	if string(response) != "pong" {
		t.Errorf("Expected pong, got %s", string(response))
	}
}
