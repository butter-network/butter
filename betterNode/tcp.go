package betterNode

import "net"

type TCPCommunication struct{}

func NewTCPCommunication() (*TCPCommunication, error) {
	return &TCPCommunication{}, nil
}

func (tc *TCPCommunication) Listen() (net.Listener, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	defer listener.Close()

	return listener, nil
}

func (tc *TCPCommunication) Connect(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
