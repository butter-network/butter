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

func PingLAN(quit chan bool, node *Node) {
	fmt.Println("in PingLAN")
	addr, err := net.ResolveUDPAddr("udp", addrGroup)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	myPingAddr = c.LocalAddr().String()
	fmt.Println("PingLAN address is: ", myPingAddr)
	myListenAddr := node.ip + ":" + node.port
	for {
		select {
		case <-quit:
			return
		default:
			//c.Write([]byte("/listening_at " + myListenAddr))
			fmt.Fprint(c, "/listening_at "+myListenAddr)
			time.Sleep(1 * time.Second)
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
	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		// do not listen to my own LAN ping
		fmt.Println("ReadFromUDP: ", src)
		fmt.Println("ReadFromUDP: ", node.knownHosts)
		srcAddrString := src.String()
		if srcAddrString != myPingAddr {
			h(src, n, b, node)
			// Stop find pinging and multicast listening
			//startUpSequenceFlag <- false
			l.Close()
			break
		}
		if len(node.knownHosts) != 0 {
			// Stop find pinging and multicast listening
			l.Close()
			break
		}
	}
}
