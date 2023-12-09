package betterNode

import "net"

type PipeCommunication struct{}

func NewPipeCommunication() (*PipeCommunication, error) {
	return &PipeCommunication{}, nil
}

func (pc *PipeCommunication) Listen() (net.Listener, error) {
	listener, err := net.Listen("unix", "/tmp/pipe.sock")
	if err != nil {
		return nil, err
	}

	defer listener.Close()

	return listener, nil
}

func (pc *PipeCommunication) Connect(addr string) (net.Conn, error) {
	conn, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
