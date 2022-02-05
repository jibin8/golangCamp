package module09

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// 指定数据包的前4个字节作为包头，里面存储的是发送的数据的长度

func Encode(message string) ([]byte, error) {
	// 读取信息的长度，转换成4个字节
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)

	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}

	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func LengthFieldHandleServer() {
	conn, err := net.Dial("tcp", ServerPort)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		msg := "test length field based frame!"
		data, err := Encode(msg)
		if err != nil {
			fmt.Println("Encode msg failed, err:", err)
			return
		}
		conn.Write(data)
	}
}

func Decode(reader *bufio.Reader) (string, error) {
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 读取真正的消息
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}

func lengthFieldHandleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Decode message failed, err:", err)
			return
		}
		fmt.Println("Received message from client, the message is:", msg)
	}
}

func LengthFieldHandleClient() {
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
		go lengthFieldHandleConn(conn)
	}
}
