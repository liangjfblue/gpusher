/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

import "time"

var (
	defaultOptions = Options{
		Address: "8888",
		Network: "tcp",
		Timeout: 5 * time.Second,
	}
)

type Options struct {
	Address         string
	Network         string
	KeepAlivePeriod time.Duration
	Timeout         time.Duration
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
