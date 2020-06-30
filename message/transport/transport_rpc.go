/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package transport

import (
	"context"
	"fmt"
	"net"
	"time"

	models2 "github.com/liangjfblue/gpusher/message/models"

	"github.com/liangjfblue/gpusher/common/utils"

	"github.com/liangjfblue/gpusher/common/discovery"

	"github.com/liangjfblue/gpusher/message/service"

	pb "github.com/liangjfblue/gpusher/proto/message/rpc/v1"

	"github.com/liangjfblue/gpusher/message/common"

	"github.com/liangjfblue/gpusher/common/codes"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"google.golang.org/grpc"
)

//暴露rpc, 提供推送接口
type rpcTransport struct {
	opts Options
}

func NewRpcTransport(opts ...Option) ITransport {
	g := new(rpcTransport)
	g.opts = defaultOptions

	for _, o := range opts {
		o(&g.opts)
	}

	return g
}

func (t *rpcTransport) Init(opts ...Option) {
	for _, o := range opts {
		o(&t.opts)
	}
}

func (t *rpcTransport) ListenServer(ctx context.Context) error {
	lis, err := net.Listen(t.opts.Network, fmt.Sprintf(":%d", t.opts.RpcPort))
	if err != nil {
		return err
	}

	go func() {
		if err = t.serve(ctx, lis); err != nil {
			log.GetLogger(common.MessageLog).Error("transport serve error, %v", err)
			panic(err)
		}
	}()

	return nil
}

func (t *rpcTransport) serve(ctx context.Context, lis net.Listener) error {
	log.GetLogger(common.MessageLog).Debug("rpc server start success, port:%d", t.opts.RpcPort)

	listener, ok := lis.(*net.TCPListener)
	if !ok {
		return codes.ErrNetworkNotSupported
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s := grpc.NewServer()

	ip, _ := utils.ExternalIP()
	etcdRegister := discovery.NewRegister(t.opts.DiscoveryAddr, -1)

	if err := etcdRegister.Register(ctx, discovery.ServiceDesc{
		ServiceName: t.opts.SrvName,
		Host:        ip,
		Port:        t.opts.RpcPort,
		TTL:         time.Second * 3,
	}); err != nil {
		return err
	}

	//注册message rpc服务
	models := models2.NewRedisModel(t.opts.RedisHost)
	if err := models.Init(); err != nil {
		return err
	}
	pb.RegisterMessageServer(s, service.NewMessageRpc(models))

	return s.Serve(listener)
}
