package traverse

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/butter-network/butter/utils"
)

// The Ambassador is a means of overcoming the NAT traversal problem. It bridges different subnetworks together by
// introducing them to each other.

type Ambassador struct {
	waitingHosts []string
	lock         sync.Mutex
}

func startAmbassador(port int16) {
	ambassador := Ambassador{}

	localIp := utils.GetOutboundIP()

	// Create listener socket
	l, err := net.Listen("tcp", localIp+":"+strconv.Itoa(int(port)))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes. (https://gobyexample.com/defer)
	defer l.Close()

	fmt.Println("Ambassador is available at ", l.Addr())
	fmt.Println("Make sure you have a Router (Port Forward) and firewall open to allow connections from other computers across the NAT")

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleAmbassadorialRequest(conn, &ambassador)
	}
}

func handleAmbassadorialRequest(conn net.Conn, ambassador *Ambassador) {
	//buf := make([]byte, 1024)
	//_, err := conn.Read(buf)
	//if err != nil {
	//	fmt.Println("Error reading:", err.Error())
	//}
	//uri, payload := parsePacket(string(buf))
	//ambassador.lock.Lock()
	//switch uri {
	//case "/get-host":
	//	// pop the first host in the queue and send it back to the node (so that they can correct directly)
	//	if len(ambassador.waitingHosts) > 0 {
	//		host := ambassador.waitingHosts[0]
	//		ambassador.waitingHosts = append(ambassador.waitingHosts[1:], "")
	//		fmt.Fprintf(conn, host)
	//	} else {
	//		fmt.Fprintf(conn, "No hosts available")
	//	}
	//case "/publish-host":
	//	ambassador.waitingHosts = append(ambassador.waitingHosts, payload)
	//}
	//ambassador.lock.Unlock()
}

// keep track of ambassador nodes in known hosts so a traversed node can ask any of it's known hosts if they either are an ambassador or know of an ambassador
