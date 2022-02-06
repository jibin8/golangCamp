package module09

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
)

type Operation uint32

type GoimInfo struct {
	Length        uint32
	HeaderLength  uint16
	Version       uint16
	Operation     Operation
	SequenceID    uint32
	Body          []byte
}

const (
	HeartBeat Operation = 2
)

func decodeLength(p *GoimInfo, r io.Reader) error {
	lengths := make([]byte, 6)
	n, err := r.Read(lengths)
	if err != nil {
		return err
	}
	if n < len(lengths) {
		return errors.New("read not enough length data")
	}
	p.Length = decodeUnit32(lengths)
	p.HeaderLength = decodeUnit16(lengths[4:])
	return nil

}

func decodeHeader(p *GoimInfo, r io.Reader) error {
	header := make([]byte, p.HeaderLength-6)
	n, err := r.Read(header)
	if err != nil {
		return err
	}
	if n < len(header) {
		return errors.New("read not enough header data")
	}
	p.Version = decodeUnit16(header)
	p.Operation = Operation(decodeUnit32(header[2:]))
	p.SequenceID = decodeUnit32(header[6:])
	return nil
}


func decodeUnit16(b []byte) uint16 {

	return binary.BigEndian.Uint16(b)

}

func decodeUnit32(b []byte) uint32 {

	return binary.BigEndian.Uint32(b)
}

func GoimDecode(r io.Reader) error {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	p := GoimInfo{}
	err := decodeLength(&p, br)
	if err != nil {
		return err
	}
	err = decodeHeader(&p, br)
	if err != nil {
		return err
	}
	p.Body = make([]byte, p.Length-uint32(p.HeaderLength))
	n, err := br.Read(p.Body)
	if err != nil {
		return err
	}
	if n < len(p.Body) {
		return errors.New("cannot read enough body data")
	}
	fmt.Printf("The goim package decoded as:\npackageLen: %d, headerLen: %d, version %d," +
		" operation: %d, sequenceId: %d, body: %v\n", p.Length, p.HeaderLength, p.Version,
		p.Operation, p.SequenceID, string(p.Body))
	return nil

}

func encode(msg string) []byte {
	headerLen := 16
	operation := HeartBeat
	version, sequence := 1, 10000
	packageLen := headerLen + len(msg)
	pkg := make([]byte, packageLen)
	binary.BigEndian.PutUint32(pkg[:4], uint32(packageLen))
	binary.BigEndian.PutUint16(pkg[4:6], uint16(headerLen))
	binary.BigEndian.PutUint16(pkg[6:8], uint16(version))
	binary.BigEndian.PutUint32(pkg[8:12], uint32(operation))
	binary.BigEndian.PutUint32(pkg[12:16], uint32(sequence))
	copy(pkg[16:], msg)
	return pkg

}

func Server() {
	conn, err := net.Dial("tcp", ServerPort)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	for i := 0; i < 3; i++ {
		msg := "test goim hhhhhhhhh dddddddd mmmmmmm"
		data := encode(msg)
		conn.Write(data)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		err := GoimDecode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Decode message failed, err:", err)
			return
		}
	}
}

func Client()  {
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
		go handleConn(conn)
	}
}