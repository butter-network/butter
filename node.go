package butter

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Block struct {
	keywords []string
	// i.e. part 1 of 5 parts
	part  int
	parts int
	data  string
}

type Node struct {
	ip               string
	port             string
	maxHostKnowledge int
	knownHosts       []string
	storage          map[string]Block
	lock             sync.Mutex
}

// NewNode creates a new node by
// - Finding the IP address of the machine
// - Assigning a port number (either user or OS defined)
// - Determining the maximum number of hosts that can be known
// - Making the known hosts slice
// - Making the storage map
func NewNode(port int) Node {
	// Determine the preferred local ip of the machine
	localIpString := GetOutboundIP().String()

	// Use defined port or default to port allocated by OS
	// port := 0

	// Determine the upper limit of the known_hosts list size based on machine resources
	maxHostKnowledge := 35

	return Node{
		ip:               localIpString,
		port:             strconv.Itoa(port),
		maxHostKnowledge: maxHostKnowledge,
		knownHosts:       make([]string, 0),
		storage:          make(map[string]Block),
	}
}

func StartNode(node *Node, clientBehaviour func(*Node), serverBehaviour func(*Node, string) string) {
	// at the same time:
	// - call out for other nodes (multicast)
	// - generate thread-pool + start listening for connections and respond to them with the prescribed listening behaviour
	// - run client behaviour as prescribed
	// TODO: Make broadcasts interact with tcpListener (maybe change the broadcast address?)
	go findPeer(node)
	//go clientBehaviour(node)
	tcpListen(node, serverBehaviour)
}

func tcpListen(node *Node, serverBehaviour func(*Node, string) string) {
	// Create listener socket
	node.lock.Lock()
	l, err := net.Listen("tcp", node.ip+":"+node.port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes. (https://gobyexample.com/defer)
	defer l.Close()

	_, port, _ := net.SplitHostPort(l.Addr().String())
	node.port = port
	node.lock.Unlock()

	fmt.Println("Listening on ", l.Addr())

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, node, serverBehaviour)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, node *Node, serverBehaviour func(*Node, string) string) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	//// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	//fmt.Println(buf)
	//// Send a response back to person contacting us.
	//conn.Write([]byte("Message received."))
	//// Close the connection when you're done with it.
	//conn.Close()
	if !routeHandler(string(buf), node) {
		// if none of the pre-defined routes match, then carry out the user defined behaviour
		serverBehaviour(node, string(buf))
	}
}

func introduceMyself(node *Node, remoteHost string) {
	// Start a tcp client connection and send them my ip and port
	c, err := net.Dial("tcp", remoteHost)
	if err != nil {
		fmt.Println(err)
		return
	}
	//c.Write([]byte"/introduction " + node.ip + ":" + node.port)
	//c.Read() // introductin confirmed
	fmt.Fprint(c, "/introduction "+node.ip+":"+node.port)
	c.Close()
}

func foundNodeHandler(src *net.UDPAddr, n int, b []byte, node *Node) {
	log.Println(n, "bytes read from", src)
	packet := string(b[:n])
	//fmt.Println(packet)
	routeHandler(packet, node)
}

// If I get a multicast that isn't myself then add it to the known hosts and stop pinging and listening
func findPeer(node *Node) {
	// Set the start-up sequence flag to true as the node is starting up and hence trying to find peers
	quit := make(chan bool, 0)

	fmt.Println("in findPeer")

	go PingLAN(quit, node)
	ListenForMulticasts(node, foundNodeHandler)
	quit <- true
	fmt.Println("I should have made a friend ", len(node.knownHosts))
}

func routeHandler(packet string, node *Node) bool {
	uri, payload := parsePacket(packet)
	switch uri {
	case "/listening_at":
		remoteHostAddress := payload
		node.lock.Lock()
		node.knownHosts = append(node.knownHosts, remoteHostAddress)
		introduceMyself(node, remoteHostAddress)
		node.lock.Unlock()
		return true
	case "/introduction": // TODO: stop ping and udp listening from here
		remoteHostAddress := payload
		node.lock.Lock()
		node.knownHosts = append(node.knownHosts, remoteHostAddress)
		node.lock.Unlock()
		return true
	}
	return false
}

func parsePacket(packet string) (string, string) {
	// get the uri by splitting the packet at the first space
	uri := strings.Split(packet, " ")[0]
	uriLength := len(uri)
	startOfPayload := uriLength + 1
	// get the payload by getting everything after the first space
	payload := packet[startOfPayload:]

	return uri, payload
}
