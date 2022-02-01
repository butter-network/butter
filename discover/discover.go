// Discovery protocol implementation for butter nodes to prevent the cold start problem (i.e. node's exist on the same
// network but are not aware of each other's existence). This package is a good example of how to build packages on top
// of butter nodes.

package discover

import (
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

func pingReceived(node *node.Node, addr []byte) []byte {
	remoteAddr, _ := utils.AddrFromJson(addr)
	node.AddKnownHost(remoteAddr)
	socketAddr := node.SocketAddr()
	nodeAddr, _ := socketAddr.ToJson()
	uri := []byte("pong/")
	_, err := utils.Request(remoteAddr, uri, nodeAddr)
	if err != nil {
		return []byte("")
	}
	return []byte("ok")
}

func pongReceived(node *node.Node, addr []byte) []byte {
	remoteAddr, err := utils.AddrFromJson(addr)
	if err != nil {
		log.Printf("pongReceived: %s", err)
		return nil
	}
	node.AddKnownHost(remoteAddr)
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
		c.Write(append(uri, socketAddress...))
		time.Sleep(1 * time.Second)

		// If I know a peer, I do not need to continue pinging the LAN
		if len(node.KnownHosts()) > 0 {
			break
		}
	}
}

func foundNode(src *net.UDPAddr, n int, b []byte, node *node.Node) {
	packet := b[:n]
	node.RouteHandler(packet)
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
		}
	}
}
