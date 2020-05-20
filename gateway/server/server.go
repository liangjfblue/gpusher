/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/liangjfblue/gpusher/gateway/config"
)

type Server struct {
	Config *config.Config
	Logger logger.ILogger
}

func NewServer(c *config.Config, l logger.ILogger) IServer {
	return &Server{
		Config: c,
		Logger: l,
	}
}

func (s *Server) Init() error {
	//初始化etcd

	//注册grpc服务

	//生成gatewayId, 向注册中心注册

	//加载缓存logic节点列表

	//初始化定时调度线程

	//初始化负载监控线程

	return nil
}

func (s *Server) Run() {
	//启动worker线程池, 处理客户端连接

	//监听链接(tcp, websocket)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)

	s.Logger.Debug("=====gateway start success=====")
	<-ch

	s.Stop()
}

func (s *Server) Stop() {
	//清理资源

	//断开etcd

	//注销grpc服务注册

	//断开监听服务
}
