/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package connect

import (
	"container/list"
	"errors"
	"sync"

	"github.com/liangjfblue/gpusher/common/codec"
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/defind"
)

var (
	MaxSubscriberPerChannel = 10000
)

var (
	ErrMaxSubscriberPerChannel = errors.New("error max subscriber per channel")
	ErrTypeConn                = errors.New("failed type Connection")
	ErrNoThisSubConn           = errors.New("no this sub Connection")
)

//ConnChannel 每个client的读写channel
type ConnChannel struct {
	mutex *sync.RWMutex
	cl    *list.List //一个key可以被多个client订阅
	num   int
}

func NewConnChannel() IChannel {
	return &ConnChannel{
		mutex: &sync.RWMutex{},
		cl:    list.New(),
		num:   0,
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
	for i := u.cl.Front(); i != nil; i = u.cl.Front().Next() {
		c := i.Value.(*Connection)
		go c.WriteMsg2Connect(key, msg)
	}
}

//创建一个客户端连接
func (u *ConnChannel) AddConn(key string, conn *Connection) (*list.Element, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	//判断当前channel分片是否达到最大conn保存数
	u.num = u.cl.Len()
	if u.num > MaxSubscriberPerChannel {
		return nil, ErrMaxSubscriberPerChannel
	}

	//连接成功首次心跳回复
	cc := codec.GetCodec(codec.Default)
	heartbeatReply, err := cc.Encode(&codec.FrameHeader{MsgType: 0x01}, nil)
	if err != nil {
		log.GetLogger(defind.GatewayLog).Error("codec Encode data err:%s", err.Error())
		return nil, err
	}

	if _, err := conn.conn.Write(heartbeatReply); err != nil {
		return nil, err
	}

	conn.HandleWriteMsg2Connect(key)

	//TODO redis保存当前网关连接数

	//client conn 加入订阅key的链表
	e := u.cl.PushFront(conn)
	u.num++

	log.GetLogger(defind.GatewayLog).Debug("user add key:%s, now sub key conn num:%d", key, u.num)

	return e, nil
}

//DelConn 删除客户端连接抽象(客户端close时调用)
func (u *ConnChannel) DelConn(key string, e *list.Element) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if e == nil {
		return
	}

	for i := u.cl.Front(); i != nil; i = u.cl.Front().Next() {
		if e.Value == i.Value {
			u.cl.Remove(e)

			cc := e.Value.(*Connection)
			close(cc.msgChan)
			u.num--
			break
		}
	}

	//去掉订阅key的对应下标的client
	log.GetLogger(defind.GatewayLog).Debug("del user conn channel key:%s, now sub key conn num:%d", key, u.num)
}

//Close 关闭所有客户端连接, 删除所有客户端抽象(server退出时主动调用)
func (u *ConnChannel) Close() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for i := u.cl.Front(); i != nil; i = u.cl.Front().Next() {
		c := i.Value.(*Connection)
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
