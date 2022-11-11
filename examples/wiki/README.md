# Butter wiki demo

This repository contains a demonstration for a wiki-style decentralised application built using the Butter framework. Two interfaces are included for the wiki, a CLI interface which you can execute with `go run cli.go` and a browser interface which you can execute using `go run browser.go` (the browser interface can be seen by navigating to [localhost:8000](http://localhost:8000/) in a browser).

Try to run the interface simultaneously, add information to one and retrieve information in the other. Then take down the original node and attempt to retrieve the information again - notice that the information persists beyond the node instance.