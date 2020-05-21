/**
 *
 * @author liangjf
 * @create on 2020/5/21
 * @version 1.0
 */
package log

import (
	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/liangjfblue/gpusher/common/logger/factory"
)

var _defaultLog logger.ILogger

func init() {
	_defaultLog = new(factory.VLogFactor).CreateLog(
		logger.Name("gpusher.log"),
		logger.Level(6),
	)
	_defaultLog.Init()
}

//RegisterLogger 暴露函数覆盖默认logger
func RegisterLogger(log logger.ILogger) {
	_defaultLog = log
	_defaultLog.Init()
}

func Trace(format string, args ...interface{}) {
	_defaultLog.Trace(format, args...)
}
func Debug(format string, args ...interface{}) {
	_defaultLog.Debug(format, args...)
}
func Info(format string, args ...interface{}) {
	_defaultLog.Info(format, args...)
}
func Warn(format string, args ...interface{}) {
	_defaultLog.Warn(format, args...)
}
func Error(format string, args ...interface{}) {
	_defaultLog.Error(format, args...)
}
