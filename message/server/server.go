/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/liangjfblue/gpusher/logic/api"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/common/transport"
	"github.com/liangjfblue/gpusher/message/common"
	"github.com/liangjfblue/gpusher/message/config"
)

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	config      *config.Config
	serviceName string

	rpcTransport transport.ITransport
}

func NewServer(c *config.Config, serviceName string) IServer {
	s := new(Server)

	s.config = c
	s.serviceName = serviceName
	return s
}

func (s *Server) Init() error {
	s.ctx, s.cancelFunc = context.WithCancel(context.TODO())

	etcdAddr := strings.Split(s.config.Server.DiscoveryAddr, ",")

	//初始化rpc客户端
	if err := api.InitMessageClientRpc(s.ctx, etcdAddr, common.MessageServiceName); err != nil {
		return err
	}

	//注册grpc服务, 暴露推送rpc接口
	s.rpcTransport = transport.NewFactoryRPCTransport(
		transport.Addr(fmt.Sprintf(":%d", s.config.Server.RpcPort)),
		transport.Network(s.config.Server.Network),
		transport.RpcPort(s.config.Server.RpcPort),
		transport.DiscoveryAddr(etcdAddr),
		transport.SrvName(s.serviceName),
	)

	return nil
}

func (s *Server) Run() {
	defer s.Stop()

	//启动rpc服务
	if err := s.rpcTransport.ListenServer(context.TODO()); err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)

	log.GetLogger(common.MessageLog).Debug("message start success")
	<-ch
}

func (s *Server) Stop() {
	log.GetLogger(common.MessageLog).Debug("message Stop clean")

	api.CloseRpcClient()
	s.cancelFunc()
}
