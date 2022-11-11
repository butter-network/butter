package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/butter-network/butter"
	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/store/pcg"
)

// WikiUser enables access to the pcg overlay API across different http request methods
type WikiUser struct {
	overlayInterface *pcg.Peer
}

// store information in the network
func (user *WikiUser) store(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("article")
	uuid := pcg.Store(user.overlayInterface, data)
	fmt.Fprintf(w, uuid)
}

// retrieve information from the network
func (user *WikiUser) retrieve(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	// Using naive retrieve for teh moment as it does a full BFS through the network - this should later be updated
	data, err := pcg.NaiveRetrieve(user.overlayInterface, strings.TrimSpace(uuid))
	if err != nil {
		fmt.Fprintf(w, "Unable to find the information on the network")
	} else {
		fmt.Fprintf(w, string(data))
	}
}

// Navigate to the interface to add a wiki entry
func addEntry(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"pages/add.html",
		"pages/base.html",
	}
	temp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}
	temp.Execute(w, nil)
}

// Navigate to the interface to find an entry
func findEntry(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"pages/find.html",
		"pages/base.html",
	}
	temp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}
	temp.Execute(w, nil)
}

// Welcom greeting page
func hello(w http.ResponseWriter, req *http.Request) {
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"pages/welcome.html",
		"pages/base.html",
	}
	temp, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
	}
	temp.Execute(w, nil)
}

func main() {
	// Crete a new node before starting the interface - this enables the user to interact with the network and in exchange provided resources
	butterNode, _ := node.NewNode(0, 512)
	fmt.Println("Node created with address:", butterNode.Address())

	// Overlay the persistent information and retrieval mechanisms
	overlay := pcg.NewPCG(butterNode, 512) // Creates a new overlay network
	pcg.AppendRetrieveBehaviour(overlay.Node())
	pcg.AppendGroupStoreBehaviour(overlay.Node())

	// Spawn the node into the network - this is blocking
	go butter.Spawn(&overlay, false, false)

	user := WikiUser{&overlay}

	// URIs
	http.HandleFunc("/", hello)
	http.HandleFunc("/add", addEntry)
	http.HandleFunc("/find", findEntry)
	http.HandleFunc("/store", user.store)
	http.HandleFunc("/retrieve", user.retrieve)

	// Start the interface server
	http.ListenAndServe(":8000", nil) // blocking
}
