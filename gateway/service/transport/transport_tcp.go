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

	"github.com/liangjfblue/gpusher/common/codes"

	"github.com/lunny/log"
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

	//Reactor模型, 为listen socket开启goroutine
	go func() {
		if err = t.serve(ctx, lis); err != nil {
			log.Errorf("transport serve error, %v", err)
		}
	}()

	return nil
}

func (t *tcpTransport) serve(ctx context.Context, lis net.Listener) error {
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

		go t.dealConn(ctx, wrapConn(conn))
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

func (t *tcpTransport) dealConn(ctx context.Context, conn *connWrapper) error {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		//frame, err := t.read(ctx, conn)
		//if err == io.EOF {
		//	return nil
		//}

		//if err != nil {
		//	return err
		//}

		//rsp, err := t.handle(ctx, frame)
		//if err != nil {
		//	return err
		//}

		////写返回结构给客户端
		//if err = t.write(ctx, conn, rsp); err != nil {
		//	return err
		//}
	}
}
