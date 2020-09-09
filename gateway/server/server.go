/**
 *
 * @author liangjf
 * @create on 2020/5/20
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
	"time"

	"github.com/liangjfblue/gpusher/common/server"

	"github.com/liangjfblue/gpusher/gateway/api"
	"github.com/liangjfblue/gpusher/gateway/service/connect"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/gateway/service/transport"

	"github.com/liangjfblue/gpusher/gateway/config"
)

type Server struct {
	config      *config.Config
	serviceName string

	transport    transport.ITransport
	rpcTransport transport.ITransport
}

func NewServer(c *config.Config, serviceName string) server.IServer {
	s := new(Server)

	s.config = c
	s.serviceName = serviceName
	return s
}

func (s *Server) Init() error {
	//初始化客户端本地缓存
	connect.InitClientChannel(s.config)

	etcdAddr := strings.Split(s.config.Server.DiscoveryAddr, ",")

	//初始化message rpc客户端
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	if err := api.InitMessageClientRpc(ctx, etcdAddr, common.MessageServiceName); err != nil {
		log.GetLogger(common.GatewayLog).Debug("gateway rpc to message err:%s", err.Error())
		return err
	}

	//TODO 初始化定时调度线程	初始化负载监控线程

	//注册grpc服务, 暴露推送rpc接口
	s.rpcTransport = transport.NewFactoryRPCTransport(
		transport.Addr(fmt.Sprintf(":%d", s.config.Server.RpcPort)),
		transport.Network(s.config.Server.Network),
		transport.RpcPort(s.config.Server.RpcPort),
		transport.DiscoveryAddr(etcdAddr),
		transport.SrvName(s.serviceName),
	)
	//选择服务器
	switch s.config.Server.Proto {
	case common.TcpProtocol:
		s.transport = transport.NewFactoryTcpTransport(
			transport.Addr(fmt.Sprintf(":%d", s.config.Server.Port)),
			transport.Network(s.config.Server.Network),
			transport.KeepAlivePeriod(time.Second*3),
		)
	case common.WsProtocol:
		s.transport = transport.NewFactoryWSTransport()
	default:
		panic("not support server type")
	}

	return nil
}

func (s *Server) Run() {
	defer s.Stop()

	//启动rpc服务
	if err := s.rpcTransport.ListenServer(context.TODO()); err != nil {
		panic(err)
	}

	//启动网关服务
	if err := s.transport.ListenServer(context.TODO()); err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)

	log.GetLogger(common.GatewayLog).Debug("gateway start success")
	<-ch
}

func (s *Server) Stop() {
	log.GetLogger(common.GatewayLog).Debug("gateway Stop clean")
	connect.GetClientChannel().Close()
	api.CloseRpcClient()
}
