package main

import (
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/defind"
	"github.com/liangjfblue/gpusher/gateway/server"

	"github.com/liangjfblue/gpusher/common/logger/factory"

	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/liangjfblue/gpusher/gateway/config"
)

func main() {
	//初始刷配置
	//win 	H:\go_home\opensource\gpusher\gateway\cmd\conf.yml
	//linux "./conf.yml"
	c := config.Init("H:\\go_home\\opensource\\gpusher\\gateway\\cmd\\conf.yml")

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
