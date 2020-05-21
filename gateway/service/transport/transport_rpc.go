/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package transport

import (
	"context"
	"net"

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
	lis, err := net.Listen(t.opts.Network, t.opts.RpcPort)
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

func (t *rpcTransport) serve(ctx context.Context, lis net.Listener) error {
	log.Debug("=====rpc server start success, port:%s=====", t.opts.RpcPort)

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

	//注册gateway rpc服务
	//pb.RegisterGreeterServer(s, &GatewayRPC{})

	return s.Serve(listener)
}

type GatewayRPC struct {
}

//New 创建新的客户端channel, 用于推送
func (g *GatewayRPC) New() {

}

//Close 关闭客户端channel
func (g *GatewayRPC) Close() {

}

//PushOne 推送给某App某用户
func (g *GatewayRPC) PushOne() {

}

//PushApp 推送给某App所有用户
func (g *GatewayRPC) PushApp() {

}

//PushAll 推送给所有App所有用户
func (g *GatewayRPC) PushAll() {

}
