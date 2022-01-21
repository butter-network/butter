# butter
> The network that spreads! 🧈

![compile_deploy_latex](https://github.com/a-shine/butter/actions/workflows/compile_deploy_latex.yml/badge.svg)

## Background

`butter` is a networking stack and library for building decentralised applications. The goal of the project is to design an 'efficient decentralised platform' that is as close as possible to being fully decentralised (no dingle controlling entity or point of failure) hence the resulting system has an unstructured peer-to-peer architecture.

Please see the full  [project documentation](https://a-shine.github.io/butter/) for more information.

**Understanding the goal 'efficient decentralised platform':**
- `butter` is efficient in the sense that it 
  - maintains persistent data consistently (despite node failure) while trying to minimise data redundancy (low degree of duplicate data)
  - takes a space-efficient approach to creating and maintaining the list known hosts per node (by determining who are the 'best' remote hosts to know on an individual node-by-node basis).
- `butter` is decentralised in the sense that it has an unstructured peer-to-peer architecture
- `butter` is a platform in the sense that it provides tooling and utilities (in the form of a library) as well as documentation and examples for developing decentralised applications

This project should help developers build decentralised systems and accommodate peers with a whole variety of hardware.

## Usage

### Preamble: communicating over the wider internet
The platform assumes a basic understanding of how to port forward (for home use). This is to deal with the pesky issue of NAT traversal when listening out for incoming requests. Port forwarding is actually very simple but may seem daunting so here is a good [guide](https://portforward.com/router.htm) on how to do it.

The library will either assign or expect the user to provide a port when creating a node. If that node needs to be accessible outside the subnetwork it will be necessary to port forward from the router to the node's assigned/chosen port.

It is worth noting that not every node port needs to be forwarded, all it takes is one node on the LAN to be accessible from the internet for the subnetwork to interact across the internet.

### Installation

1. From within your Go project, run the following command:
   ```bash
   go get github.com/a-shine/butter
   ```
   This will download the library and configure your `go.mod` file.
2. Import the package into your project source:
   ```go
   import "github.com/a-shine/butter"
   ```

### Examples
#### Non-persistent data (Chat application)
```go
// TODO: Add commented example code
```
#### Persistent data (Wiki application)
```go
// TODO: Add commented example code
````

Take a look at more examples in the [examples/](./examples) directory.

## Working on the platform

### Development

- See the project board
- Raise an issue or pull request

**Explaining the project structure**:
Making a peer to peer system can be broken down into 5 main components
1. Defining the node behaviour and maintaining the network by managing known hosts (per node) (`node/`)
2. Discovery (`discover/`)
3. NAT traversal (`traverse/`)
4. (Persistent) information storage (`store/`)
5. (Persistent) information retrieval (`retrieve/`)
These are reflected in the project package structure

### Testing

When testing on a development machine, it may become necessary to test the behaviour across several IP addresses, this can be achieved by running nodes from different containers.