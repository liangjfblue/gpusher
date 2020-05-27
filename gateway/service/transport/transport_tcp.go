/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"time"

	"github.com/liangjfblue/gpusher/gateway/proto"

	"github.com/liangjfblue/gpusher/common/codec"

	"github.com/liangjfblue/gpusher/gateway/defind"

	"github.com/liangjfblue/gpusher/gateway/service"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/common/codes"
)

var (
	HeartbeatReply             = []byte("ok")
	ErrConnReqPayloadInfoReply = []byte("conn req payload info err")
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
			log.GetLogger(defind.GatewayLog).Error("transport serve error, %v", err)
		}
	}()

	return nil
}

func (t *tcpTransport) serve(ctx context.Context, lis net.Listener) error {
	log.GetLogger(defind.GatewayLog).Debug("=====tcp server start success, port:%s=====", t.opts.Address)

	listener, ok := lis.(*net.TCPListener)
	if !ok {
		return codes.ErrNetworkNotSupported
	}

	var delayTmp time.Duration
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := listener.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if delayTmp == 0 {
					delayTmp = 5 * time.Millisecond
				} else {
					delayTmp *= 2
				}
				if max := 1 * time.Second; delayTmp > max {
					delayTmp = max
				}
				time.Sleep(delayTmp)
				continue
			}
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
	defer func() {
		conn.Close()
		if err := recover(); err != nil {
			log.GetLogger(defind.GatewayLog).Error("client error")
		}
	}()

	select {
	case <-ctx.Done():
		log.GetLogger(defind.GatewayLog).Error(ctx.Err().Error())
		return
	default:
	}

	addr := conn.RemoteAddr().String()
	log.GetLogger(defind.GatewayLog).Debug("new conn coming, addr:%s", addr)

	//读取连接client信息
	framer, err := t.read(conn)
	if err != nil {
		if err == io.EOF {
			log.GetLogger(defind.GatewayLog).Warn("client conn close err:%s", err.Error())
		} else {
			log.GetLogger(defind.GatewayLog).Error("read err:%s", err.Error())
		}
		return
	}

	//decode data
	cc := codec.GetCodec(codec.Default)
	connReq, err := cc.Decode(framer)
	if err != nil {
		log.GetLogger(defind.GatewayLog).Error("decode payload err:%s", err.Error())
		return
	}

	//检验appId key token参数
	var connPayload proto.ConnPayload
	if err = json.Unmarshal(connReq, &connPayload); err != nil {
		log.GetLogger(defind.GatewayLog).Error("connect gateway payload err:%s", err.Error())
		return
	}

	if connPayload.AppId == 0 || connPayload.Key == "" || connPayload.Token == "" {
		log.GetLogger(defind.GatewayLog).Error("conn req payload info error")
		if _, err := conn.Conn.Write(ErrConnReqPayloadInfoReply); err != nil {
			log.GetLogger(defind.GatewayLog).Error("conn req payload info err, err:%s", err.Error())
			return
		}
	}

	appId, key, token := connPayload.AppId, connPayload.Key, connPayload.Token
	log.GetLogger(defind.GatewayLog).Debug("appId:%d, key:%s, token:%s", appId, key, token)

	//TODO 优化心跳检测间隔, 检查heartbeat的间隔

	//TODO 检验token

	//创建一个Connection结构代替原始conn, 并等待channel的推送消息
	userConn, err := service.GetClientChannel().Get(appId, key, true)
	if err != nil {
		log.GetLogger(defind.GatewayLog).Error("get userConn err:%s", err.Error())
		return
	}

	//把key对应的connection加入对应appChannel, 创建一个goroutine负责写推送消息给客户端
	connection := service.NewConnect(conn, defind.TcpProtocol, codec.GetVersion(framer))
	index, err := userConn.AddConn(key, connection)
	if err != nil {
		log.GetLogger(defind.GatewayLog).Error("add user conn channel err:%s", err.Error())
		return
	}
	defer func() {
		if err := userConn.DelConn(key, index); err != nil {
			log.GetLogger(defind.GatewayLog).Error("del user conn channel err:%s", err.Error())
			return
		}
	}()

	for {
		//read heartbeat
		framer, err := t.read(conn)
		if err != nil {
			if err == io.EOF {
				log.GetLogger(defind.GatewayLog).Warn("client conn close err:%s", err.Error())
			} else {
				log.GetLogger(defind.GatewayLog).Error("for read err:%s", err.Error())
			}
			break
		}

		if codec.IsHeartBeatMsg(framer) {
			cc := codec.GetCodec(codec.Default)
			resp, err := cc.Encode(&codec.FrameHeader{MsgType: 0x01}, nil)
			if err != nil {
				log.GetLogger(defind.GatewayLog).Error("codec Encode data err:%s", err.Error())
				return
			}

			if _, err := conn.Conn.Write(resp); err != nil {
				log.GetLogger(defind.GatewayLog).Error("conn write HeartbeatReply, err:%s", err.Error())
				return
			}
		}
	}
}

func (t *tcpTransport) read(conn *connWrapper) ([]byte, error) {
	return conn.framer.ReadFramer()
}
