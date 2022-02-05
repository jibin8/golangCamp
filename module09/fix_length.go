package module09

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const (
	FixLength = 1024
	ServerPort    = "127.0.0.1:3001"
)

func fixLengthHandleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		var msg = make([]byte, FixLength)
		n, err := reader.Read(msg)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("read message failed, err:", err)
			return
		}
		fmt.Println("the message is:", string(msg[:n]))
	}
}

func FixLengthClient() {
	listen, err := net.Listen("tcp", ServerPort)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go fixLengthHandleConn(conn)
	}
}

func FixLengthServer() {
	conn, err := net.Dial("tcp", ServerPort)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	for i := 0; i < 50; i++ {
		msg := []byte("test fix length")
		fixLength := FixLength - len(msg)
		msg = append(msg, make([]byte, fixLength)...)
		conn.Write(msg)
	}
}
