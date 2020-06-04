/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package connect

import "container/list"

//IChannel 客户端分配channel通信
//不直接用conn维护客户端连接, 而是用channel, 是为了和通信协议解耦, 支持任意通信协议(tcp, ws, udp...)
type IChannel interface {
	//AddToken 客户端连接添加token权限
	AddToken(string, string) error
	//CheckToken 校验客户端连接token权限
	CheckToken(string, string) error
	//PushMsg 推送消息
	PushMsg(int, string, []byte) error
	//创建一个客户端连接
	AddConn(int, string, *Connection) (*list.Element, error)
	//删除一个客户端连接
	DelConn(int, string, *list.Element)
	//Close 关闭客户channel
	Close() error
}
