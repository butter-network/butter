package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const SocketAddressSize int = 6 // bytes

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

func AddrFromJson(addressInJson []byte) (SocketAddr, error) {
	socketAddress := SocketAddr{}
	err := json.Unmarshal(addressInJson, &socketAddress)
	if err != nil {
		err := errors.New("unable to convert the json address to a socket address")
		return socketAddress, err
	}
	return socketAddress, nil
}