// Discovery protocol implementation for butter nodes to prevent the cold start problem (i.e. node's exist on the same
// network but are not aware of each other's existence). This package is a good example of how to build packages on top
// of butter nodes.

package discover

import (
	"encoding/json"
	"github.com/butter-network/butter/node"
	"github.com/butter-network/butter/utils"
	"log"
	"net"
	"time"
)

const pingInterval = 1
const checkConnected = 5

const (
	pingRoute       = "ping/"
	pongRoute       = "pong/"
	aliveRoute      = "alive/"
	addrGroup       = "224.0.0.1:9999" // Multicast address for all nodes
	maxDatagramSize = 8192
)

func alive(overlay node.Overlay, addr []byte) []byte {
	return []byte("ok")
}

func pingReceived(overlay node.Overlay, addr []byte) []byte {
	remoteAddr, _ := utils.AddrFromJson(addr)
	overlay.Node().AddKnownHost(remoteAddr)
	addrs := make([]utils.SocketAddr, 0)
	addrs = append(addrs, overlay.Node().SocketAddr())
	for addr := range overlay.Node().KnownHosts() {
		addrs = append(addrs, addr)
	}
	knownHostsJson, _ := json.Marshal(addrs)
	uri := []byte("pong/")
	_, err := utils.Request(remoteAddr, uri, knownHostsJson)
	if err != nil {
		return []byte("")
	}
	return []byte("ok")
}

func pongReceived(overlay node.Overlay, addr []byte) []byte {
	addrs := make([]utils.SocketAddr, 0)
	err := json.Unmarshal(addr, &addrs)
	//remoteAddr, err := utils.AddrFromJson(addr)
	if err != nil {
		log.Printf("pongReceived: %s", err)
		return nil
	}
	for _, addr := range addrs {
		if addr != overlay.Node().SocketAddr() {
			overlay.Node().AddKnownHost(addr)
		}
	}
	//overlay.Node().AddKnownHost(remoteAddr)
	return []byte("/successful-introduction/")
}

func Discover(overlay node.Overlay) {
	overlay.Node().RegisterServerBehaviour(pingRoute, pingReceived)
	overlay.Node().RegisterServerBehaviour(pongRoute, pongReceived)

	overlay.Node().RegisterServerBehaviour(aliveRoute, alive)

	go ListenForMulticasts(overlay)
	go checkImStillConnected(overlay)
	PingLAN(overlay)

}

var myPingAddr string

func PingLAN(overlay node.Overlay) {
	addr, err := net.ResolveUDPAddr("udp", addrGroup)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	myPingAddr = c.LocalAddr().String()
	uri := []byte(pingRoute)
	socketAddr := overlay.Node().SocketAddr()
	socketAddress, _ := socketAddr.ToJson()
	for {
		c.Write(append(uri, socketAddress...))
		time.Sleep(pingInterval * time.Second)

		// If I know a peer, I do not need to continue pinging the LAN
		if len(overlay.Node().KnownHosts()) > 0 {
			break
		}
	}
}

func foundNode(src *net.UDPAddr, n int, b []byte, overlay node.Overlay) {
	packet := b[:n]
	overlay.Node().RouteHandler(packet, overlay)
}

func ListenForMulticasts(overlay node.Overlay) {
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
			foundNode(src, n, b, overlay)
		}
	}
}

func checkImStillConnected(overlay node.Overlay) {
	for {
		time.Sleep(checkConnected * time.Second)
		for addr := range overlay.Node().KnownHosts() {
			_, err := utils.Request(addr, []byte("alive/"), nil)
			if err != nil {
				overlay.Node().RemoveKnownHost(addr) // This is also being done in the known host update function - should I have it twice?
			}
		}
		if len(overlay.Node().KnownHosts()) == 0 {
			//fmt.Println("No known hosts, restarting")
			PingLAN(overlay)
		} else {
			//fmt.Println("Yay, I'm connected!")
		}
	}
}
