/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package logger

import "time"

var DefaultOptions = Options{
	Name:     "gpusher",
	Dir:      ".",
	Level:    6, //TraceLevel
	Interval: 1,
	Host:     "127.0.0.1",
	TimeOut:  3 * time.Second,
}
