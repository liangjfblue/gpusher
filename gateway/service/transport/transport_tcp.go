/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

import (
	"context"
	"io"
	"net"

	"github.com/liangjfblue/gpusher/common/codec"

	"github.com/liangjfblue/gpusher/gateway/defind"

	"github.com/liangjfblue/gpusher/gateway/service"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/common/codes"
)

var (
	HeartbeatReply = []byte("ok")
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
	select {
	case <-ctx.Done():
		log.Error(ctx.Err().Error())
		return
	default:
	}

	defer conn.Close()

	addr := conn.RemoteAddr().String()
	log.Debug("new conn coming, addr:%s", addr)

	//TODO 读取一帧, 解析得到key, token
	framer, err := t.read(conn)
	if err == io.EOF {
		log.Error("read compeleted")
		return
	}

	if err != nil {
		log.Error("first read data framer err:%s", err.Error())
		return
	}

	//if codec.IsHeartBeatMsg(framer) {
	//
	//}

	appId, key, token, version, heartbeat := 0, "", "", "", 0
	log.Debug(key, token, version, heartbeat)

	//TODO 优化心跳检测间隔, 检查heartbeat的间隔

	//TODO 检验token

	//创建一个Connection结构代替原始conn, 并等待channel的推送消息
	userConn, err := service.GetUserChannel().Get(appId, key, true)
	if err != nil {
		log.Error("get userConn err:%s", err.Error())
		return
	}

	//把key对应的connection加入对应appChannel
	connection := service.NewConnect(conn, defind.TcpProtocol, codec.GetVersion(framer))
	userConn.AddConn(key, connection)
	defer userConn.DelConn(key)

	for {
		framer, err := t.read(conn)
		if err == io.EOF {
			log.Error("read compeleted")
			return
		}

		if codec.IsHeartBeatMsg(framer) {
			cc := codec.GetCodec(codec.Default)
			resp, err := cc.Encode(&codec.FrameHeader{MsgType: 0x01}, nil)
			if err != nil {
				log.Error("codec Encode data err:%s", err.Error())
				return
			}

			if _, err := conn.framer.Write(resp); err != nil {
				log.Error("conn write HeartbeatReply, err:%s", err.Error())
				return
			}
		}
	}
}

func (s *tcpTransport) read(conn *connWrapper) ([]byte, error) {
	return conn.framer.ReadFramer()
}
