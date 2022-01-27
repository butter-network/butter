package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
)

const EOF byte = 26

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

func createConnections(remoteHost SocketAddr) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteHost.ToString())
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Read(conn net.Conn) ([]byte, error) {
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer
	for {
		b, err := reader.ReadByte()
		fmt.Printf("%d,", b)
		if err != nil {
			return nil, err
		}
		if b == EOF {
			fmt.Println(" - I'm here!!!")
			break
		}
		buffer.WriteByte(b)
	}
	fmt.Println("done reading") // BUG: Not getting to here sometimes
	return buffer.Bytes(), nil
}

func Write(conn net.Conn, packet []byte) error {
	writer := bufio.NewWriter(conn)
	appended := append(packet, EOF)
	_, err := writer.Write(appended)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	fmt.Println("done writing")
	return nil
}

func Request(remoteHost SocketAddr, route []byte, payload []byte) ([]byte, error) {
	conn, err := createConnections(remoteHost)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	packet := append(route, payload...)
	packet = append(packet, EOF)

	err = Write(conn, packet)
	if err != nil {
		return nil, err
	}

	response, err := Read(conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}
