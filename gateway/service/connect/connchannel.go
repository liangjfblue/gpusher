/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package connect

import (
	"errors"
	"sync"

	"github.com/liangjfblue/gpusher/common/codec"
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/defind"
)

var (
	MaxSubscriberPerChannel = 32
)

var (
	ErrMaxSubscriberPerChannel = errors.New("error max subscriber per channel")
	ErrTypeConn                = errors.New("failed type Connection")
	ErrNoThisSubConn           = errors.New("no this sub Connection")
)

//ConnChannel 每个client的读写channel
type ConnChannel struct {
	mutex *sync.RWMutex
	cl    []interface{} //一个key可以被多个client订阅
}

func NewConnChannel() IChannel {
	return &ConnChannel{
		mutex: &sync.RWMutex{},
	}
}

//AddToken 客户端连接添加token权限
func (u *ConnChannel) AddToken(string, string) error {
	return nil
}

//CheckToken 校验客户端连接token是否超时
func (u *ConnChannel) CheckToken(string, string) error {
	return nil
}

//PushMsg 写推送消息到通道
func (u *ConnChannel) PushMsg(key string, msg []byte) error {
	//TODO 私有消息
	u.write(key, msg)
	return nil
}

//Write 写推送消息到通道
func (u *ConnChannel) Write(key string, msg []byte) error {
	u.write(key, msg)
	return nil
}

func (u *ConnChannel) write(key string, msg []byte) {
	//发送给多有订阅key的client
	for _, v := range u.cl {
		c := v.(*Connection)
		go c.WriteMsg2Connect(key, msg)
	}
}

//创建一个客户端连接
func (u *ConnChannel) AddConn(key string, conn *Connection) (int, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	//判断当前channel分片是否达到最大conn保存数
	if len(u.cl) > MaxSubscriberPerChannel {
		return 0, ErrMaxSubscriberPerChannel
	}

	//连接成功首次心跳回复
	cc := codec.GetCodec(codec.Default)
	heartbeatReply, err := cc.Encode(&codec.FrameHeader{MsgType: 0x01}, nil)
	if err != nil {
		log.GetLogger(defind.GatewayLog).Error("codec Encode data err:%s", err.Error())
		return 0, err
	}

	if _, err := conn.conn.Write(heartbeatReply); err != nil {
		return 0, err
	}

	conn.HandleWriteMsg2Connect(key)

	//TODO redis保存当前网关连接数

	//client conn 加入订阅key的链表
	u.cl = append(u.cl, conn)
	connSubKeyIndex := len(u.cl) - 1

	log.GetLogger(defind.GatewayLog).Debug("user add key:%s, sub index:%d, now sub key conn num:%d", key, connSubKeyIndex, len(u.cl))

	return connSubKeyIndex, nil
}

//DelConn 删除客户端连接抽象(客户端close时调用)
func (u *ConnChannel) DelConn(key string, index int) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	//TODO bug, index是会变化的
	if len(u.cl) <= index {
		return ErrNoThisSubConn
	}

	cc, ok := u.cl[index].(*Connection)
	if !ok {
		return ErrTypeConn
	}

	close(cc.msgChan)

	//去掉订阅key的对应下标的client
	u.cl = append(u.cl[:index], u.cl[index+1:]...)
	log.GetLogger(defind.GatewayLog).Debug("del user conn channel key:%s, index:%d, now sub key conn num:%d", key, index, len(u.cl))
	return nil
}

//Close 关闭所有客户端连接, 删除所有客户端抽象(server退出时主动调用)
func (u *ConnChannel) Close() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for _, v := range u.cl {
		c := v.(*Connection)
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	u.cl = u.cl[:0]

	return nil
}
