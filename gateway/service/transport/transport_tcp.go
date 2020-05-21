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

	"github.com/liangjfblue/gpusher/gateway/defind"

	"github.com/liangjfblue/gpusher/gateway/service"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/common/codes"
)

type tcpTransport struct {
	opts Options
}

func NewTcpTransport(opts ...Option) ITransport {
	t := new(tcpTransport)
	t.opts = defaultOptions

	for _, o := range opts {
		o(&t.opts)
	}

	return t
}

func (t *tcpTransport) Init(opts ...Option) {
	for _, o := range opts {
		o(&t.opts)
	}
}

func (t *tcpTransport) ListenServer(ctx context.Context) error {
	lis, err := net.Listen(t.opts.Network, t.opts.Address)
	if err != nil {
		return err
	}

	go func() {
		if err = t.serve(ctx, lis); err != nil {
			log.Error("transport serve error, %v", err)
		}
	}()

	return nil
}

func (t *tcpTransport) serve(ctx context.Context, lis net.Listener) error {
	log.Debug("=====tcp server start success, port:%s=====", t.opts.Address)

	listener, ok := lis.(*net.TCPListener)
	if !ok {
		return codes.ErrNetworkNotSupported
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}

		conn, err = t.setConn(conn)
		if err != nil {
			return err
		}

		go t.dealTCPConn(ctx, wrapConn(conn))
	}
}

func (t *tcpTransport) setConn(conn *net.TCPConn) (*net.TCPConn, error) {
	if err := conn.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if t.opts.KeepAlivePeriod > 0 {
		if err := conn.SetKeepAlivePeriod(t.opts.KeepAlivePeriod); err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func (t *tcpTransport) dealTCPConn(ctx context.Context, conn *connWrapper) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	log.Debug("new conn coming, addr:%s", addr)

	//TODO 读取一帧, 解析得到key, token
	key, token, version, heartbeat := "", "", "", 0
	log.Debug(key, token, version, heartbeat)

	//TODO 检查heartbeat的间隔

	//TODO 检验token

	//创建一个Connection结构代替原始conn, 并等待channel的推送消息
	connection := service.NewConnect(conn, defind.TcpProtocol, version)
	connection.HandleWriteMsg(key)

	//把key对应的connection加入对应appChannel

	begin := time.Now().UnixNano()
	end := begin + int64(time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Error(ctx.Err().Error())
		default:
		}

		//间隔性检查heartbeat有效性, 超过时间
		if end-begin >= int64(time.Second) {
			if err := conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(heartbeat))); err != nil {
				log.Error("<%s> key:%s, error: %v", addr, key, err)
				break
			}
			begin = end
		}
	}
}
