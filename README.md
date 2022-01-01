# Butter
> The network that spreads! 🧈

![compile_deploy_latex](https://github.com/a-shine/butter/actions/workflows/compile_deploy_latex.yml/badge.svg)

Full project documentation - [https://a-shine.github.io/butter/](https://a-shine.github.io/butter/)

## Brief introduction to the project

Butter is a project with the goal of building an 'efficient decentralised platform'.

- The resulting platform takes the form of a library that should enable the easy development of decentralised applications by abstracting away most of the decentralised behaviour.
- The platform is efficient as it strives to minimise redundant data replication while still maintaining data stability despite node failure. In addition, the platform takes a space-efficient approach to creating and maintaining the list known hosts per node (determine who is best to know on a individual node-by-node basis). This should reduce the barrier to accessibility and accommodate peers with a whole variety of hardware.

The project's unique design approach is to think about events (meeting a new peer, interacting with a peer, conversing with a peer, spreading information) on the network in a social way, drawing inspiration from modelling human interaction in a social context. The information on the network should emulate the way humans naturally meet, communicate and dissipate information. A social gathering (e.g. a party or a casual meeting of friends) is arguably a good model of decentralised communication and hence can be drawn upon as inspiration.

The platform assumes an (albeit fairly rudimentary) understanding of port forwarding - this is a very simple skill which may seem daunting to new users but it really is quite simple - here is a good guide...

It is worth noting that the library expects you to provide a port when instantiating a node and if you would like to make that particular node accessible outside the sub-network you are expected to port forward from the router to the chosen port.

## Using the library

### Installing

1. From within your Go project, run the following command:
   ```bash
   go get github.com/a-shine/butter
   ```
   This will download the library and configure your `go.mod` file.
2. Simply import the package into your project:
   ```go
   import "github.com/a-shine/butter"
   ```

### Example use with non-persistent data - decentralised chat

```go
// ADD CODE HERE WHEN COMPLETE
```

### Example use with persistent data (i.e. using the data storage functionality) - simplistic decentralised encyclopedia

```go
// ADD CODE HERE WHEN COMPLETE
```

Take a look at more examples in the [examples dir](./examples).

## Working on the platform

### Developing

- see the project board
- raise an issue to contribute

### Testing

the platform on a network locally

We want to be testing with different IPs (at least on the local network). We can use containers to have several IPs and test the system.

To open an interactive shell
```bash
docker-compose run node bash
```