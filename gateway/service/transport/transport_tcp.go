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
	"sync"
	"time"

	"github.com/liangjfblue/gpusher/gateway/service/message"

	"github.com/liangjfblue/gpusher/gateway/service/connect"

	"github.com/liangjfblue/gpusher/gateway/proto"

	"github.com/liangjfblue/gpusher/common/codec"

	"github.com/liangjfblue/gpusher/gateway/common"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/common/codes"
)

var (
	HeartbeatReply             = []byte("ok")
	ErrConnReqPayloadInfoReply = []byte("conn req payload info err")
)

type tcpTransport struct {
	opts Options
	pool sync.Pool
}

func NewTcpTransport(opts ...Option) ITransport {
	t := new(tcpTransport)
	t.opts = defaultOptions
	t.pool.New = func() interface{} {
		return t.allocateWrapConn()
	}

	for _, o := range opts {
		o(&t.opts)
	}

	return t
}

func (t *tcpTransport) allocateWrapConn() *connWrapper {
	return &connWrapper{}
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
			log.GetLogger(common.GatewayLog).Error("transport serve error, %v", err)
		}
	}()

	return nil
}

func (t *tcpTransport) serve(ctx context.Context, lis net.Listener) error {
	log.GetLogger(common.GatewayLog).Debug("tcp server start success, port%s", t.opts.Address)

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

		//use connWrapper object pool
		connWrapper := t.pool.Get().(*connWrapper)
		connWrapper.setWrapConn(conn)

		go t.ioHandle(ctx, wrapConn(conn))

		t.pool.Put(connWrapper)
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

//ioHandle 处理用户io事件
func (t *tcpTransport) ioHandle(ctx context.Context, conn *connWrapper) {
	defer func() {
		_ = conn.Close()
		if err := recover(); err != nil {
			log.GetLogger(common.GatewayLog).Error("client error: %s", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.GetLogger(common.GatewayLog).Error(ctx.Err().Error())
		return
	default:
	}

	addr := conn.RemoteAddr().String()
	log.GetLogger(common.GatewayLog).Debug("new conn coming, addr:%s", addr)

	connPayload, err := t.loginConn(conn)
	if err != nil {
		log.GetLogger(common.GatewayLog).Error(ctx.Err().Error())
		return
	}

	//TODO 检验token

	//创建一个Connection结构代替原始conn, 并等待channel的推送消息
	userConn, err := connect.GetClientChannel().Get(connPayload.AppId, connPayload.UUID, true)
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("get userConn err:%s", err.Error())
		return
	}

	//把key对应的connection加入对应appChannel, 创建一个goroutine负责写推送消息给客户端
	connection := connect.NewConnect(conn, common.TcpProtocol)
	e, err := userConn.AddConn(connPayload.AppId, connPayload.UUID, connection)
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("add user conn channel err:%s", err.Error())
		return
	}
	defer userConn.DelConn(connPayload.AppId, connPayload.UUID, e)

	for {
		//read heartbeat
		framer, err := t.read(conn)
		if err != nil {
			if err == io.EOF {
				err = nil
				log.GetLogger(common.GatewayLog).Warn("client conn close")
			} else {
				log.GetLogger(common.GatewayLog).Error("for read err:%s", err.Error())
			}
			goto EXIT
		}

		switch codec.GetMsgType(framer) {
		case codec.HeartbeatMsg:
			cc := codec.GetCodec(codec.Default)
			resp, err := cc.Encode(&codec.FrameHeader{MsgType: codec.HeartbeatMsg}, nil)
			if err != nil {
				log.GetLogger(common.GatewayLog).Error("codec Encode data err:%s", err.Error())
				goto EXIT
			}

			if _, err := conn.Conn.Write(resp); err != nil {
				log.GetLogger(common.GatewayLog).Error("conn write HeartbeatReply, err:%s", err.Error())
				goto EXIT
			}

			//TODO rpc to message 续期redis的路由, 防止gateway还保留旧的路由映射
			if err := message.ExpireGatewayUUID(connPayload.UUID); err != nil {
				log.GetLogger(common.GatewayLog).Error("ExpireGatewayUUID err:%s", err.Error())
			}
		case codec.MsgAckMsg:
			//TODO 客户端确认消费, rpc到logic确认
		case codec.SyncOfflineMsg:
		//TODO 获取离线消息
		default:
			log.GetLogger(common.GatewayLog).Error("the msg type not support")
		}
	}

EXIT:
	//rpc message删除路由
	if err := message.DeleteGatewayUUID(connPayload.UUID); err != nil {
		log.GetLogger(common.GatewayLog).Error("DeleteGatewayUUID err:%s", err.Error())
	}
}

//read 读完整一帧数据
func (t *tcpTransport) read(conn *connWrapper) ([]byte, error) {
	return conn.framer.ReadFramer()
}

//loginConn 用户连接后必须发送特定一帧数据表示连接初始化
func (t *tcpTransport) loginConn(conn *connWrapper) (*proto.ConnPayload, error) {
	//读取连接client信息
	framer, err := t.read(conn)
	if err != nil {
		if err == io.EOF {
			log.GetLogger(common.GatewayLog).Warn("client conn close err:%s", err.Error())
		} else {
			log.GetLogger(common.GatewayLog).Error("read err:%s", err.Error())
		}
		return nil, err
	}

	//decode data
	cc := codec.GetCodec(codec.Default)
	connReq, err := cc.Decode(framer)
	if err != nil {
		log.GetLogger(common.GatewayLog).Error("decode payload err:%s", err.Error())
		return nil, err
	}

	//检验appId uuid key token参数
	var connPayload proto.ConnPayload
	if err = json.Unmarshal(connReq, &connPayload); err != nil {
		log.GetLogger(common.GatewayLog).Error("connect gateway payload err:%s", err.Error())
		return nil, err
	}

	if connPayload.AppId == 0 || connPayload.UUID == "" || connPayload.Key == "" || connPayload.Token == "" {
		log.GetLogger(common.GatewayLog).Error("conn req payload info error")
		if _, err := conn.Conn.Write(ErrConnReqPayloadInfoReply); err != nil {
			log.GetLogger(common.GatewayLog).Error("conn req payload info err, err:%s", err.Error())
			return nil, err
		}
	}
	return &connPayload, nil
}
