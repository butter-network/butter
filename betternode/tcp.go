package betternode

import (
	"bufio"
	"bytes"
	"log"
	"net"
)

type TCP struct {
	addr     string
	port     int
	listener net.Listener
}

func (n *TCP) Listen() {
	for {
		_, err := n.listener.Accept()
		if err != nil {
			// Avoid fatal errors at all costs - we want to maximise node availability
			log.Println("Node is unable to accept incoming connections due to: ", err.Error())
			continue // forces next iteration of the loop skipping any code in between
		}

		// Pass connection to request handler in a new goroutine - allows a node to handle multiple connections at once
		//go node.HandleRequest(conn, overlay)
	}
}

func (n *TCP) Request(commInterface CommunicationInterface, route []byte, payload []byte) ([]byte, error) {
	tcpInterface := commInterface.(*TCP)
	conn, err := createConnections(tcpInterface.addr, tcpInterface.port)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	packet := append(route, payload...)
	packet = append(packet, EOF)

	err = Write(&conn, packet)
	if err != nil {
		return nil, err
	}

	response, err := Read(&conn)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func createConnections(address string, port int) (net.Conn, error) {
	socketAddr := address + ":" + string(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", socketAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
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
