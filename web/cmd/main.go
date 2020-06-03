package main

import (
	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/liangjfblue/gpusher/common/logger/factory"
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/web/config"
	"github.com/liangjfblue/gpusher/web/server"
)

func main() {
	c := config.Init("./conf.yml")

	l := new(factory.VLogFactor).CreateLog(
		logger.Name(c.Log.Name),
		logger.Level(c.Log.Level),
	)
	log.RegisterLogger(common.WebLog, l)

	s := server.NewServer(c)
	if err := s.Init(); err != nil {
		panic(err)
	}

	s.Run()
}
