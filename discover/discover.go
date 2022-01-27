package discover

import (
	"fmt"
	"github.com/a-shine/butter/node"
	"github.com/a-shine/butter/utils"
	"log"
	"net"
	"time"
)

const (
	pingRoute       = "ping/"
	pongRoute       = "pong/"
	addrGroup       = "224.0.0.1:9999"
	maxDatagramSize = 8192
)

// All routes have node and payload as parameters and return a response.

func pingReceived(node *node.Node, addr []byte) []byte {
	log.Printf("ping received from %s", addr)
	remoteAddr, _ := utils.AddrFromJson(addr)
	node.AddNewKnownHost(remoteAddr)
	socketAddr := node.SocketAddr()
	nodeAddr, _ := socketAddr.ToJson()
	uri := []byte("pong/")
	request, err := utils.Request(remoteAddr, uri, nodeAddr)
	if err != nil {
		log.Printf("ping request failed: %s", err)
		return []byte("")
	}
	fmt.Println("asking for pong: ", request)
	return []byte("ok")
}

func pongReceived(node *node.Node, addr []byte) []byte {
	log.Printf("pong received from %s", addr)
	remoteAddr, err := utils.AddrFromJson(addr)
	if err != nil {
		log.Printf("pongReceived: %s", err)
		return nil
	}
	node.AddNewKnownHost(remoteAddr)
	return []byte("/successful-introduction/")
}

func Discover(node *node.Node) {
	node.RegisterRoute(pingRoute, pingReceived)
	node.RegisterRoute(pongRoute, pongReceived)

	go ListenForMulticasts(node)
	PingLAN(node)

}

var myPingAddr string

func PingLAN(node *node.Node) {
	addr, err := net.ResolveUDPAddr("udp", addrGroup)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	myPingAddr = c.LocalAddr().String()
	uri := []byte(pingRoute)
	socketAddr := node.SocketAddr()
	socketAddress, _ := socketAddr.ToJson()
	for {
		fmt.Println("I'm pinging...")
		c.Write(append(uri, socketAddress...))
		time.Sleep(1 * time.Second)

		// If I know a peer, I do not need to continue pinging the LAN
		if len(node.KnownHosts()) > 0 {
			fmt.Println("I know a peer, so I am done pinging the LAN")
			break
		}
	}
}

func foundNode(src *net.UDPAddr, n int, b []byte, node *node.Node) {
	packet := b[:n]
	node.NewRouteHandler(packet)
}

func ListenForMulticasts(node *node.Node) {
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

		srcAddrString := src.String()
		if srcAddrString != myPingAddr {
			foundNode(src, n, b, node)
			fmt.Println("Known peers: ", node.KnownHosts())
		}
	}
}
