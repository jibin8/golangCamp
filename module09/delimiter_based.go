package module09

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

const Tag = '\n'

func delimiterBasedHandleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadSlice(Tag)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("read message failed, err:", err)
			return
		}
		fmt.Println("the message is:", string(msg))
	}
}

func DelimiterBasedClient() {
	listen, err := net.Listen("tcp", ServerPort)
	if err != nil {
		fmt.Println("Listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept failed, err:", err)
			continue
		}
		go delimiterBasedHandleConn(conn)
	}
}

func DelimiterBasedServer() {
	conn, err := net.Dial("tcp", ServerPort)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		msg := []byte("test delimiter!\n")
		conn.Write(msg)
	}
}
