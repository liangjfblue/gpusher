/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package logger

import "time"

type Options struct {
	Name     string
	Dir      string
	Level    int
	Interval int

	Host    string
	TimeOut time.Duration
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Dir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

func Level(level int) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func Interval(interval int) Option {
	return func(o *Options) {
		o.Interval = interval
	}
}

func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}
func TimeOut(timeOut time.Duration) Option {
	return func(o *Options) {
		o.TimeOut = timeOut
	}
}
