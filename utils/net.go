package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const SocketAddressSize = 6 // in bytes

// SocketAddr has a size 4 bytes (IP) + 2 bytes (port) = 6 bytes
type SocketAddr struct {
	Ip   net.IP
	Port uint16
}

func (s *SocketAddr) ToString() string {
	return fmt.Sprintf("%s:%d", s.Ip.String(), s.Port)
}

func (s *SocketAddr) ToJson() ([]byte, error) {
	e, err := json.Marshal(s)
	return e, err
}

func FromJson(addressInJson []byte) (SocketAddr, error) {
	socketAddress := SocketAddr{}
	err := json.Unmarshal(addressInJson, &socketAddress)
	return socketAddress, err
}

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
