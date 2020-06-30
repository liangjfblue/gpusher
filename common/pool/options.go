/**
 *
 * @author liangjf
 * @create on 2020/6/8
 * @version 1.0
 */
package pool

var (
	defaultOptions = Options{
		initCap: 3,
		maxCap:  10,
		maxIdle: 1,
	}
)

type Options struct {
	initCap int
	maxCap  int
	maxIdle int
	builder Builder
}

func WithInitCap(initCap int) Option {
	return func(o *Options) {
		o.initCap = initCap
	}
}

func WithMaxCap(maxCap int) Option {
	return func(o *Options) {
		o.maxCap = maxCap
	}
}

func WithMaxIdle(maxIdle int) Option {
	return func(o *Options) {
		o.maxIdle = maxIdle
	}
}

func WithBuilder(builder Builder) Option {
	return func(o *Options) {
		o.builder = builder
	}
}
