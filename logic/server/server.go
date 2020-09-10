/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package server

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/liangjfblue/gpusher/common/server"
	"github.com/liangjfblue/gpusher/logic/api"
	"github.com/liangjfblue/gpusher/logic/models"

	"github.com/liangjfblue/gpusher/logic/service"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/logic/common"

	"github.com/liangjfblue/gpusher/logic/config"
)

type Server struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	config      *config.Config
	serviceName string
}

func NewServer(c *config.Config, serviceName string) server.IServer {
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

	//初始化gateway rpc所有客户端
	if err := service.InitGatewayRpcClient(etcdAddr); err != nil {
		return err
	}

	//初始化message rpc客户端
	if err := api.InitMessageClientRpc(ctx, etcdAddr, common.MessageServiceName); err != nil {
		log.GetLogger(common.LogicLog).Debug("logic rpc to message err:%s", err.Error())
		return err
	}

	redisAddr := strings.Split(s.config.Redis.Host, ",")
	if err := models.InitRedisModel(redisAddr); err != nil {
		log.GetLogger(common.LogicLog).Debug("init redis err:%s", err.Error())
		return err
	}

	ctx, cancel = context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	kafkaAddr := strings.Split(s.config.Kafka.BrokerAddrs, ",")
	go func() {
		if err := service.InitKafkaConsumer(ctx, kafkaAddr); err != nil {
			return
		}
	}()

	return nil
}

func (s *Server) Run() {
	defer s.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)

	log.GetLogger(common.LogicLog).Debug("logic start success")
	<-ch
}

func (s *Server) Stop() {
	log.GetLogger(common.LogicLog).Debug("logic Stop clean")

	service.StopKafkaConsumer()
	service.CLoseRpcClient()
	s.cancelFunc()
	api.CloseRpcClient()
}
