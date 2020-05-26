package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/liangjfblue/gpusher/common/codec"
)

func main() {
	fmt.Println("client launch")

	tcpAddr, err := net.ResolveTCPAddr("tcp", "172.16.7.17:8881")
	if err != nil {
		fmt.Println("Resolve TCPAddr error", err)
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)

	if err != nil {
		fmt.Println("connect server error", err)
	}

	cc := codec.GetCodec(codec.Default)
	r, err := cc.Encode(&codec.FrameHeader{
		Magic:        codec.Magic,
		Version:      codec.Version,
		MsgType:      0x01,
		ReqType:      0x0,
		CompressType: 0x0,
		StreamID:     1,
		Reserved:     0,
		Length:       uint32(len("hello")),
	}, []byte("hello"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(r))

	if _, err := conn.Write(r); err != nil {
		log.Fatal(err)
	}
	//go recv(conn)
	time.Sleep(2 * time.Second)

	_ = conn.Close()
}

//func recv(conn net.Conn) {
//	buffer := make([]byte, 1024)
//	n, err := conn.Read(buffer)
//	if err == nil {
//		fmt.Println("read message from server:" + string(buffer[:n]))
//		fmt.Println("Message len:", n)
//	}
//}
