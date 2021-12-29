package butter

import (
	"log"
	"net"
	"strings"
)

// GetOutboundIP gets the preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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
