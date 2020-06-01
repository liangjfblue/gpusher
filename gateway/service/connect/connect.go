/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package connect

import (
	"net"

	"github.com/liangjfblue/gpusher/common/codec"

	"github.com/liangjfblue/gpusher/gateway/defind"

	"github.com/liangjfblue/gpusher/common/logger/log"
)

//Connection 客户端连接抽象层
type Connection struct {
	conn    net.Conn
	proto   string
	msgChan chan []byte
}

//NewConnect 创建客户端连接抽象
func NewConnect(conn net.Conn, proto string) *Connection {
	return &Connection{
		conn:    conn,
		proto:   proto,
		msgChan: make(chan []byte),
	}
}

//HandleWriteMsg2Connect 客户端推送消息通道监听和写客户端
func (c *Connection) HandleWriteMsg2Connect(uuid string) {
	go func() {
		var (
			n   int
			err error
		)
		for msg := range c.msgChan {
			switch c.proto {
			case defind.TcpProtocol:
				//tcp自定义协议
				var resp []byte
				cc := codec.GetCodec(codec.Default)
				resp, err = cc.Encode(&codec.FrameHeader{MsgType: 0x01}, msg)
				if err != nil {
					log.GetLogger(defind.GatewayLog).Error("codec Encode data err:%s", err.Error())
					return
				}
				n, err = c.conn.Write(resp)
			case defind.WsProtocol:
				n, err = c.conn.Write(msg)
			default:
				log.GetLogger(defind.GatewayLog).Error("not support proto type")
			}
			log.GetLogger(defind.GatewayLog).Debug("uuid:%s, write msg n:%d", uuid, n)

			if err != nil {
				log.GetLogger(defind.GatewayLog).Error("uuid:%s, write msg err:%s", uuid, err.Error())
			}
		}
		log.GetLogger(defind.GatewayLog).Debug("uuid:%s, conn goroutine closed", uuid)
	}()
}

//WriteMsg2Connect 对外暴露, 用于推送消息到chan的转换
func (c *Connection) WriteMsg2Connect(uuid string, msg []byte) {
	select {
	case c.msgChan <- msg:
		log.GetLogger(defind.GatewayLog).Debug("uuid:%s write msg:%s", uuid, string(msg))
	default:
		c.conn.Close()
		log.GetLogger(defind.GatewayLog).Error("uuid:%s write msg:%s; close conn", uuid, string(msg))
	}
}
