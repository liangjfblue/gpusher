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

	"github.com/liangjfblue/gpusher/gateway/common"

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
func (c *Connection) HandleWriteMsg2Connect(appId int, uuid string) {
	go func() {
		var (
			n   int
			err error
		)
		for msg := range c.msgChan {
			switch c.proto {
			case common.TcpProtocol:
				//tcp自定义协议
				var resp []byte
				cc := codec.GetCodec(codec.Default)
				resp, err = cc.Encode(
					&codec.FrameHeader{
						MsgType:  codec.GeneralMsg,
						StreamID: 1, //TODO 消息序号生成器 redis
					},
					msg,
				)
				if err != nil {
					log.GetLogger(common.GatewayLog).Error("codec Encode data err:%s", err.Error())
					return
				}
				n, err = c.conn.Write(resp)
			case common.WsProtocol:
				n, err = c.conn.Write(msg)
			default:
				log.GetLogger(common.GatewayLog).Error("not support proto type")
			}
			log.GetLogger(common.GatewayLog).Debug("appId:%d, uuid:%s, write msg n:%d", appId, uuid, n)

			if err != nil {
				log.GetLogger(common.GatewayLog).Error("appId:%d, uuid:%s, write msg err:%s", appId, uuid, err.Error())
			}
		}
		log.GetLogger(common.GatewayLog).Debug("appId:%d, uuid:%s, conn goroutine closed", appId, uuid)
	}()
}

//WriteMsg2Connect 对外暴露, 用于推送消息到chan的转换
func (c *Connection) WriteMsg2Connect(appId int, uuid string, msg []byte) {
	select {
	case c.msgChan <- msg:
		log.GetLogger(common.GatewayLog).Debug("appId:%d,uuid:%s write msg:%s", appId, uuid, string(msg))
	default:
		c.conn.Close()
		log.GetLogger(common.GatewayLog).Error("appId:%d,uuid:%s write msg:%s; close conn", appId, uuid, string(msg))
	}
}
