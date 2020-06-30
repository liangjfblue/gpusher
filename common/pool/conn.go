/**
 *
 * @author liangjf
 * @create on 2020/6/8
 * @version 1.0
 */
package pool

import (
	"errors"
	"net"
	"sync"
)

var (
	ErrNetHadClosed = errors.New("net had closed")
)

//Conn 原生net.Conn包装器
type poolConn struct {
	sync.RWMutex
	net.Conn
	unUseFlag bool
	pc        *poolChannel
}

func (c *poolConn) Close() error {
	c.RLock()
	defer c.RUnlock()

	if c.unUseFlag {
		if c.Conn != nil {
			return c.Conn.Close()
		}
	}

	//put to poolChannel
	c.pc.put(c.Conn)
	return c.Conn.Close()
}

func (c *poolConn) Read(buf []byte) (int, error) {
	c.RLock()
	defer c.RUnlock()

	if c.unUseFlag {
		return 0, ErrNetHadClosed
	}

	n, err := c.Conn.Read(buf)
	if err != nil {
		c.SetUnUseFlag()
		c.Conn.Close()
	}

	return n, err
}

func (c *poolConn) Write(buf []byte) (int, error) {
	c.RLock()
	defer c.RUnlock()

	if c.unUseFlag {
		return 0, ErrNetHadClosed
	}

	n, err := c.Conn.Write(buf)
	if err != nil {
		c.SetUnUseFlag()
		c.Conn.Close()
	}

	return n, err
}

func (c *poolConn) SetUnUseFlag() {
	c.Lock()
	defer c.Unlock()

	c.unUseFlag = true
}

func (p *poolChannel) WrapConn(conn net.Conn) net.Conn {
	c := &poolConn{
		pc:        p,
		unUseFlag: false,
	}
	c.Conn = conn
	return c
}
