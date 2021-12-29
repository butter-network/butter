package butter

import (
	"encoding/json"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"net"
	"os"
	"strconv"
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
	go findPeer(node)
	go clientBehaviour(node)
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
	//quit := make(chan bool, 0)

	fmt.Println("in findPeer")

	go ListenForMulticasts(node, foundNodeHandler) // TODO: This should actually always be listening for nodes that want to join the network
	PingLAN(node)                                  // This should stop once it has found a peer

	//quit <- true
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

func GetKnownHosts(node *Node) []string {
	node.lock.Lock()
	defer node.lock.Unlock()
	return node.knownHosts
}

func Send(remoteHost string, message string) {
	// Start a tcp client connection and send them my ip and port
	c, err := net.Dial("tcp", remoteHost)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(c, message)
	reply := make([]byte, 1024)
	c.Read(reply)
	fmt.Println("Reply:", string(reply))
	c.Close()
}

// NaiveRetrieve High level entrypoint for searching for a specific piece of information on the network
// look if I have the information else look at the most likely known host to get to that information
// one query per piece of information (one-to-one) hence the query has to be unique i.e i.d.
func NaiveRetrieve(node *Node, query string) string {
	// do I have this information, if so return it
	// else BFS (pass the query on to all known hosts (partial view)
	node.lock.Lock()
	defer node.lock.Unlock()
	if val, ok := node.storage[query]; ok {
		return val.data
	} else {
		return bfs(node, query)
	}
}

func bfs(node *Node, query string) string {
	// Initialise an empty queue
	queue := make([]string, 0)
	// Add all my known hosts to the queue
	for _, host := range node.knownHosts {
		queue = append(queue, host)
	}
	for len(queue) > 0 {
		// Pop the first element from the queue
		host := queue[0]
		queue = queue[1:]
		// Start a connection to the host
		c, err := net.Dial("tcp", host)
		if err != nil {
			fmt.Println(err)
			return "Error connecting to host"
		}
		c.Close()
		// Ask host if he has data
		fmt.Fprint(c, "/remote-retrieve "+query)
		// Receive response
		reply := make([]byte, 1024)
		c.Read(reply)
		uri, payload := parsePacket(string(reply))
		// If the returned packet is success + the data then return it
		// else add the known hosts of the remote node to the end of the queue
		if uri == "/success" {
			return payload
		} else {
			fmt.Fprint(c, "/get-remote-known-hosts"+query)
			c.Read(reply)
			// convert json list of known hosts into a slice of strings
			remoteHosts := make([]string, 0)
			err = json.Unmarshal(reply, &remoteHosts)
			if err != nil {
				fmt.Println(err)
				return "Error decoding json"
			}
			// add the remote hosts to the end of the queue
			queue = append(queue, remoteHosts...)
		}
		return "Information is not on the network"
	}
	return "This should not happen"
}

// NaiveStore stores information on the network naively by simply placing it on the local node. It generate a UUIS for
// the information and creates an information block and return information uuid
func NaiveStore(node *Node, keywords []string, information string) string {
	node.lock.Lock()
	// Generate UUID
	u, _ := uuid.NewV4()
	node.storage[u.String()] = Block{
		keywords: keywords,
		part:     0,
		parts:    0,
		data:     information,
	}
	node.lock.Unlock()
	return u.String()
}
