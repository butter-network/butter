package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const SocketAddressSize int = 6 // Bytes

// SocketAddr has a size 4 bytes (IPv4) + 2 bytes (port) = 6 bytes
type SocketAddr struct {
	Ip   string
	Port uint16
}

func (s *SocketAddr) IsEmpty() bool {
	return s.Ip == "" || s.Port == 0
}

type SocketAddrSlice []SocketAddr

func (s *SocketAddr) ToString() string {
	return fmt.Sprintf("%s:%d", s.Ip, s.Port)
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

func AddrFromString(address string) (SocketAddr, error) {
	socketAddress := SocketAddr{}
	addressParts := strings.Split(address, ":")
	if len(addressParts) != 2 {
		err := errors.New("invalid address format")
		return socketAddress, err
	}

	// check that IP is valid
	ip := net.ParseIP(addressParts[0])
	if ip == nil {
		err := errors.New("invalid ip address")
		return socketAddress, err
	}

	port, err := strconv.Atoi(addressParts[1])
	if err != nil {
		err := errors.New("invalid port")
		return socketAddress, err
	}
	socketAddress.Ip = addressParts[0]
	socketAddress.Port = uint16(port)
	return socketAddress, nil
}

func (s *SocketAddrSlice) ToJson() ([]byte, error) {
	e, err := json.Marshal(s)
	return e, err
}

func AddrSliceFromJson(addrJson []byte) (SocketAddrSlice, error) {
	var s SocketAddrSlice
	err := json.Unmarshal(addrJson, &s)
	return s, err
}
