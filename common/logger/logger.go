/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package logger

type ILogger interface {
	Init(...Option)
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type Option func(*Options)
