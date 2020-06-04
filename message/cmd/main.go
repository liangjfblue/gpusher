package main

import (
	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/liangjfblue/gpusher/common/logger/factory"
	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/message/common"
	"github.com/liangjfblue/gpusher/message/config"
	"github.com/liangjfblue/gpusher/message/server"
)

//go:generate protoc -I ../proto/rpc/v1 --go_out=plugins=grpc:../proto/rpc/v1 ../proto/rpc/v1/api.proto

func main() {
	c := config.Init("./conf.yml")

	l := new(factory.VLogFactor).CreateLog(
		logger.Name(c.Log.Name),
		logger.Level(c.Log.Level),
	)
	log.RegisterLogger(common.MessageLog, l)

	s := server.NewServer(c, common.MessageServiceName)
	if err := s.Init(); err != nil {
		panic(err)
	}

	s.Run()
}
