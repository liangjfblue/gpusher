/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package service

import (
	"net"

	"github.com/liangjfblue/gpusher/common/codec"

	"github.com/liangjfblue/gpusher/gateway/defind"

	"github.com/liangjfblue/gpusher/common/logger/log"
)

type Connection struct {
	Conn    net.Conn
	Proto   string
	Version uint32
	MsgChan chan []byte
}

func NewConnect(conn net.Conn, proto string, version uint32) *Connection {
	return &Connection{
		Conn:    conn,
		Proto:   proto,
		Version: version,
		MsgChan: make(chan []byte),
	}
}

func (c *Connection) HandleWriteMsg(key string) {
	go func() {
		var (
			n   int
			err error
		)
		for msg := range c.MsgChan {
			switch c.Proto {
			case defind.TcpProtocol:
				//tcp自定义协议
				var resp []byte
				cc := codec.GetCodec(codec.Default)
				resp, err = cc.Encode(&codec.FrameHeader{MsgType: 0x01}, msg)
				if err != nil {
					log.GetLogger(defind.GatewayLog).Error("codec Encode data err:%s", err.Error())
					return
				}
				n, err = c.Conn.Write(resp)
			case defind.WsProtocol:
				n, err = c.Conn.Write(msg)
			default:
				log.GetLogger(defind.GatewayLog).Error("not support proto type")
			}
			log.GetLogger(defind.GatewayLog).Debug("key:%s, write msg n:%d", key, n)

			if err != nil {
				log.GetLogger(defind.GatewayLog).Error("key:%s, write msg err:%s", key, err.Error())
			}
		}
		log.GetLogger(defind.GatewayLog).Debug("key:%s, conn goroutine closed", key)
	}()
}

func (c *Connection) WriteMsg(key string, msg []byte) {
	select {
	case c.MsgChan <- msg:
		log.GetLogger(defind.GatewayLog).Debug("key:%s write msg:%s", key, string(msg))
	default:
		c.Conn.Close()
		log.GetLogger(defind.GatewayLog).Error("key:%s write msg:%s; close conn", key, string(msg))
	}
}
