/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

import (
	"time"
)

var (
	defaultOptions = Options{
		Address:       "8888",
		Network:       "tcp",
		Timeout:       5 * time.Second,
		RpcPort:       9990,
		DiscoveryAddr: []string{"172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"},
		SrvName:       "none",
	}
)

type Options struct {
	Address         string
	Network         string
	KeepAlivePeriod time.Duration
	Timeout         time.Duration
	RpcPort         int
	DiscoveryAddr   []string
	SrvName         string
}

func Addr(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func Network(network string) Option {
	return func(o *Options) {
		o.Network = network
	}
}

func KeepAlivePeriod(keepAlivePeriod time.Duration) Option {
	return func(o *Options) {
		o.KeepAlivePeriod = keepAlivePeriod
	}
}

func Timeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
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
