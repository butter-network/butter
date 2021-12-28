package main

import (
	"encoding/gob"
	"fmt"
	"github.com/a-shine/butter"
	"net"
)

func main() {
	var input string
	go butter.StartServer()
	fmt.Scanln(&input)
	go client()

	fmt.Scanln(&input)
}
func client() {
	// connect to the server

	// wait for 3 sec

	c, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		fmt.Println(err)
		return
	}
	// send the message
	msg := "Hello World"
	fmt.Println("Sending", msg)
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}
