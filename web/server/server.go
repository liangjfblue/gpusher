package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/liangjfblue/gpusher/common/push"
	"github.com/liangjfblue/gpusher/web/service"

	"github.com/liangjfblue/gpusher/web/config"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/web/router"
)

type Server struct {
	config *config.Config
	Router *router.Router
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
		Router: router.NewRouter(),
	}
}

func (s *Server) Init() error {
	//消息队列
	addr := strings.Split(s.config.Kafka.BrokerAddrs, ",")
	q := push.NewKafkaSender(addr, false)
	if err := q.Init(); err != nil {
		return err
	}

	//推送者
	service.RegisterPush("kafka", service.NewDefaultPush(q))
	go service.GetPush("kafka").Run()

	s.Router.Init()
	return nil
}

func (s *Server) Run() {
	defer func() {
		log.GetLogger(common.WebLog).Debug("web close, clean and close something")
	}()

	log.GetLogger(common.WebLog).Debug("web server Run, port:%d", s.config.Server.Port)

	log.GetLogger(common.WebLog).Error(http.ListenAndServe(fmt.Sprintf(":%d", s.config.Server.Port), s.Router.G).Error())
}

func (s *Server) Stop() {
	log.GetLogger(common.WebLog).Debug("web Stop clean")
}
