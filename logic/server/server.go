/**
 *
 * @author liangjf
 * @create on 2020/6/3
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

	"github.com/liangjfblue/gpusher/logic/service"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/logic/common"

	"github.com/liangjfblue/gpusher/common/transport"
	"github.com/liangjfblue/gpusher/logic/config"
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
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	if err := service.InitGatewayRpcClient(etcdAddr); err != nil {
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

	ctx, cancel = context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	kafkaAddr := strings.Split(s.config.Kafka.BrokerAddrs, ",")
	if err := service.InitKafkaConsumer(ctx, kafkaAddr); err != nil {
		return err
	}

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

	log.GetLogger(common.LogicLog).Debug("logic start success")
	<-ch
}

func (s *Server) Stop() {
	log.GetLogger(common.LogicLog).Debug("logic Stop clean")

	s.cancelFunc()
	service.StopKafkaConsumer()
	service.CLoseRpcClient()
}
