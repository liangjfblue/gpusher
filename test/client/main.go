package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/liangjfblue/gpusher/common/codec"
)

func main() {
	fmt.Println("client test")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "172.16.7.17:8881")
	if err != nil {
		log.Fatal("Resolve TCPAddr error", err)
	}

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatal("connect server error", err)
	}

	type ConnPayload struct {
		AppId int    `json:"appId"`
		Key   string `json:"key"`
		Token string `json:"token"`
	}

	connReq := ConnPayload{
		AppId: 1000,
		Key:   "liangjf",
		Token: "test",
	}
	jConnReq, err := json.Marshal(connReq)
	if err != nil {
		log.Fatal(err)
	}

	//encode heartbeat appId key token
	cc := codec.GetCodec(codec.Default)
	r, err := cc.Encode(
		&codec.FrameHeader{
			Magic:        codec.Magic,
			Version:      codec.Version,
			MsgType:      0x01,
			ReqType:      0x0,
			CompressType: 0x0,
			StreamID:     1,
			Reserved:     0,
		},
		jConnReq,
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Write(r); err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("read len:%d, message:%v\n", n, buffer[:n])

	time.Sleep(2 * time.Second)

	go func() {
		//encode heartbeat data
		cc := codec.GetCodec(codec.Default)
		r, err := cc.Encode(&codec.FrameHeader{
			Magic:        codec.Magic,
			Version:      codec.Version,
			MsgType:      0x00,
			ReqType:      0x0,
			CompressType: 0x0,
			StreamID:     1,
			Reserved:     0,
		}, []byte("hello"))
		if err != nil {
			log.Fatal(err)
		}

		if _, err := conn.Write(r); err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("read len:%d, message:%v\n", n, buffer[:n])
	}()

	time.Sleep(3 * time.Second)

	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}
}
