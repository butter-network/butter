package betternode

import (
	"bufio"
	"bytes"
	mock_conn "github.com/jordwest/mock-conn"
	"net"
)

type Pipe struct {
	conn mock_conn.Conn
}

func (p *Pipe) Listen() {
	listener := bufio.NewReader(p.conn.Server)
	for {
		var buffer bytes.Buffer
		for {
			b, err := listener.ReadByte()
			if err != nil {
				break
			}
			if b == EOF {
				break
			}
			buffer.WriteByte(b)
		}
		// Handle connection
	}
}

func (p *Pipe) Request(commInterface CommunicationInterface, route []byte, payload []byte) ([]byte, error) {
	pipeInterface := commInterface.(*Pipe)
	return nil, nil
}

func Read(conn *net.Conn) ([]byte, error) {
	reader := bufio.NewReader(*conn)
	var buffer bytes.Buffer
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == EOF {
			break
		}
		buffer.WriteByte(b)
	}
	return buffer.Bytes(), nil
}

func Write(conn *net.Conn, packet []byte) error {
	writer := bufio.NewWriter(*conn)
	appended := append(packet, EOF)
	_, err := writer.Write(appended)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
