/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

var (
	defaultOptions = Options{
		Addr:          "8888",
		Network:       "tcp",
		RpcPort:       9990,
		DiscoveryAddr: []string{"172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"},
		SrvName:       "none",
		RedisHost:     []string{"172.16.7.16:8001", "172.16.7.16:8002", "172.16.7.16:8003"},
	}
)

type Options struct {
	Addr          string
	Network       string
	RpcPort       int
	DiscoveryAddr []string
	SrvName       string
	RedisHost     []string
}

func Addr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func Network(network string) Option {
	return func(o *Options) {
		o.Network = network
	}
}

func RpcPort(rpcPort int) Option {
	return func(o *Options) {
		o.RpcPort = rpcPort
	}
}

func DiscoveryAddr(DiscoveryAddr []string) Option {
	return func(o *Options) {
		o.DiscoveryAddr = DiscoveryAddr
	}
}

func SrvName(SrvName string) Option {
	return func(o *Options) {
		o.SrvName = SrvName
	}
}

func RedisHost(redisHost []string) Option {
	return func(o *Options) {
		o.RedisHost = redisHost
	}
}
