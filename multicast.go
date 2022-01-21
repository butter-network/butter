package butter

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	addrGroup       = "224.0.0.1:9999"
	maxDatagramSize = 8192
)

var myPingAddr string

func PingLAN(node *Node) {
	addr, err := net.ResolveUDPAddr("udp", addrGroup)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	myPingAddr = c.LocalAddr().String()
	fmt.Println("Pinging out for peers at ", myPingAddr)
	//myListenAddr := node.ip + ":" + node.port
	//myListenAddr := node.socketAddr.ToString()
	//for {
	//	select {
	//	case <-quit:
	//		return
	//	default:
	//		//c.Write([]byte("/listening_at " + myListenAddr))
	//		fmt.Fprint(c, "/listening_at "+myListenAddr)
	//		time.Sleep(1 * time.Second)
	//	}
	//}
	for {
		fmt.Println("I'm pinging...")
		uri := []byte("/listening_at ")
		socketAddress, _ := node.socketAddr.ToJson()
		c.Write(append(uri, socketAddress...))
		time.Sleep(1 * time.Second)

		// If I know a peer, I do not need to continue pinging the LAN
		if len(node.knownHosts) > 0 {
			fmt.Println("I know a peer, so I am done pinging the LAN")
			break
		}
	}
}

// How to not listen to my own pings?

func ListenForMulticasts(node *Node, h func(*net.UDPAddr, int, []byte, *Node)) {
	addr, err := net.ResolveUDPAddr("udp", addrGroup)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	defer l.Close()
	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		// do not listen to my own LAN ping
		//fmt.Println("ReadFromUDP: ", src)
		//fmt.Println("ReadFromUDP: ", node.knownHosts)
		srcAddrString := src.String()
		if srcAddrString != myPingAddr {
			h(src, n, b, node)
			fmt.Println("Known peers: ", node.knownHosts)
			// Stop find pinging and multicast listening
			//startUpSequenceFlag <- false
			//l.Close()
			//break
		}
		//if len(node.knownHosts) != 0 {
		//	// Stop find pinging and multicast listening
		//	l.Close()
		//	break
		//}
	}
}
