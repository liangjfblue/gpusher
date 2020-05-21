/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

import (
	"context"
	"net"
	"time"
)

type ITransport interface {
	Init(...Option)
	ListenServer(context.Context) error
}

type Option func(*Options)

type connWrapper struct {
	net.Conn
	CurTime int64
	//framer Framer
}

func wrapConn(rawConn net.Conn) *connWrapper {
	return &connWrapper{
		Conn:    rawConn,
		CurTime: time.Now().Unix(),
		//framer: NewFramer(),
	}
}
