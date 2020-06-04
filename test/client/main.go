package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/liangjfblue/gpusher/common/message"

	"github.com/liangjfblue/gpusher/common/codec"
)

func main() {
	fmt.Println("client test")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "172.16.7.16:8881")
	if err != nil {
		log.Fatal("Resolve TCPAddr error", err)
	}

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatal("connect server error", err)
	}
	defer conn.Close()

	connect(conn)

	time.Sleep(100 * time.Millisecond)

	go sendHeartbeat(conn)

	go func() {
		framer := message.NewFramer(conn)

		cc := codec.GetCodec(codec.Default)
		for {
			d, err := framer.ReadFramer()
			if err != nil {
				log.Fatal(err)
			}

			if codec.IsHeartBeatMsg(d) {
				fmt.Printf("read heartbeat message:%v\n", d[:codec.FrameHeadLen])
				continue
			}

			//TODO do something what do you want to do
			data, err := cc.Decode(d)
			fmt.Printf("read push message:%v\n", string(data))
		}
	}()

	select {}
}

func connect(conn net.Conn) {
	type ConnPayload struct {
		AppId int    `json:"appId"`
		UUID  string `json:"uuid"`
		Key   string `json:"key"`
		Token string `json:"token"`
	}

	connReq := ConnPayload{
		AppId: 1000,
		UUID:  "liangjf",
		Key:   "app_gpusher",
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
}

func sendHeartbeat(conn net.Conn) {
	//encode heartbeat data
	cc := codec.GetCodec(codec.Default)
	r, err := cc.Encode(&codec.FrameHeader{
		Magic:        codec.Magic,
		Version:      codec.Version,
		MsgType:      0x01, //0x1: heartbeat
		ReqType:      0x0,
		CompressType: 0x0,
		StreamID:     1,
		Reserved:     0,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			if _, err := conn.Write(r); err != nil {
				log.Fatal(err)
			}
		}
	}
}
