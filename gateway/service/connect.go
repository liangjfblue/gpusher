/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package service

import (
	"net"

	"github.com/liangjfblue/gpusher/gateway/defind"

	"github.com/liangjfblue/gpusher/common/logger/log"
)

type Connection struct {
	Conn    net.Conn
	Proto   string
	Version string
	MsgChan chan []byte
}

func NewConnect(conn net.Conn, proto string, version string) *Connection {
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
				n, err = c.Conn.Write(msg)
			case defind.WsProtocol:
				//tcp自定义协议 封包
				n, err = c.Conn.Write(msg)
			default:
				log.Error("not support proto type")
			}
			log.Debug("key:%s, write msg n:%d", key, n)

			if err != nil {
				log.Error("key:%s, write msg err:%s", key, err.Error())
			}
		}
		log.Debug("key:%s, conn goroutine closed", key)
	}()
}

func (c *Connection) WriteMsg(key string, msg []byte) {
	select {
	case c.MsgChan <- msg:
		log.Debug("key:%s write msg:%s", key, string(msg))
	default:
		c.Conn.Close()
		log.Error("key:%s write msg:%s; close conn", key, string(msg))
	}
}
