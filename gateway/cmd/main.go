package main

import (
	"flag"

	"github.com/liangjfblue/gpusher/common/logger/log"
	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/gateway/config"
	"github.com/liangjfblue/gpusher/gateway/server"

	"github.com/liangjfblue/gpusher/common/logger/factory"

	"github.com/liangjfblue/gpusher/common/logger"
)

var (
	ServiceName = "gpusher.gateway"
	httpPort    int
	rpcPort     int
)

func init() {
	flag.IntVar(&httpPort, "hp", 8881, "http port")
	flag.IntVar(&rpcPort, "rp", 7771, "rpc port")
}

//go:generate protoc -I ../proto/rpc/v1 --go_out=plugins=grpc:../proto/rpc/v1 ../proto/rpc/v1/gateway.proto

func main() {
	flag.Parse()

	c := config.Init("./conf.yml")
	if httpPort > 0 {
		c.Server.Port = httpPort
	}
	if rpcPort > 0 {
		c.Server.RpcPort = rpcPort
	}

	l := new(factory.VLogFactor).CreateLog(
		logger.Name(c.Log.Name),
		logger.Level(c.Log.Level),
	)
	log.RegisterLogger(common.GatewayLog, l)

	s := server.NewServer(c, ServiceName)
	if err := s.Init(); err != nil {
		panic(err)
	}

	s.Run()
}
