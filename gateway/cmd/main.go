package main

import (
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/config"
	"github.com/liangjfblue/gpusher/gateway/defind"
	"github.com/liangjfblue/gpusher/gateway/server"

	"github.com/liangjfblue/gpusher/common/logger/factory"

	"github.com/liangjfblue/gpusher/common/logger"
)

func main() {
	//初始刷配置
	configT := "./conf.yml"
	//configT := "H:\\go_home\\opensource\\gpusher\\gateway\\cmd\\conf.yml"
	c := config.Init(configT)

	//初始化日志
	l := new(factory.VLogFactor).CreateLog(
		logger.Name(c.Log.Name),
		logger.Level(c.Log.Level),
	)
	log.RegisterLogger(defind.GatewayLog, l)

	s := server.NewServer(c)
	if err := s.Init(); err != nil {
		panic(err)
	}

	s.Run()
}
