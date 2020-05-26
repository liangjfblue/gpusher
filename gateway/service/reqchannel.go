/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package service

import (
	"container/list"
	"sync"
)

type UserChannel struct {
	mutex *sync.RWMutex
	cl    list.List
}

func NewUserChannel() IChannel {
	return &UserChannel{
		mutex: &sync.RWMutex{},
	}
}

//AddToken 客户端连接添加token权限
func (u *UserChannel) AddToken(string, string) error {
	return nil
}

//CheckToken 校验客户端连接token权限
func (u *UserChannel) CheckToken(string, string) error {
	return nil
}

//PushMsg 推送消息
func (u *UserChannel) PushMsg(string, []byte) error {
	return nil
}

//Write 写返回结果给客户端
func (u *UserChannel) Write(string, []byte) error {
	return nil
}

//创建一个客户端连接
func (u *UserChannel) AddConn(key string, conn *Connection) {

	conn.HandleWriteMsg(key)
	return
}

//删除一个客户端连接
func (u *UserChannel) DelConn(key string) {
	return
}

//Close 关闭客户channel
func (u *UserChannel) Close() error {
	return nil
}
