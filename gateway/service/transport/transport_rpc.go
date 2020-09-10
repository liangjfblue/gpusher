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

	"github.com/liangjfblue/gpusher/common/utils"

	"github.com/liangjfblue/gpusher/common/discovery"

	"github.com/liangjfblue/gpusher/gateway/service/rpc"

	pb "github.com/liangjfblue/gpusher/proto/gateway/rpc/v1"

	"github.com/liangjfblue/gpusher/gateway/common"

	"github.com/liangjfblue/gpusher/common/codes"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"google.golang.org/grpc"
)

//StartRPCServer 注册rpc服务,提供接口用于推送相关;
//创建message服务客户端(gateway不处理推送路有关系)

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
			log.GetLogger(common.GatewayLog).Error("transport serve error, %v", err)
		}
	}()

	return nil
}

func (t *rpcTransport) serve(ctx context.Context, lis net.Listener) error {
	log.GetLogger(common.GatewayLog).Debug("rpc server start success, port:%d", t.opts.RpcPort)

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

	//注册gateway rpc服务
	pb.RegisterGatewayServer(s, &rpc.GatewayRPC{})

	return s.Serve(listener)
}
