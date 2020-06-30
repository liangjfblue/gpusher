/**
 *
 * @author liangjf
 * @create on 2020/6/8
 * @version 1.0
 */
package pool

import "net"

//IPool 连接池管理接口
type IPool interface {
	Get() (net.Conn, error)
	Close()
}

type Option func(*Options)
