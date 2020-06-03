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

var (
	_logMap    map[string]logger.ILogger
	DefaultLog = new(factory.VLogFactor).CreateLog(
		logger.Name("gpusher.log"),
		logger.Level(6),
	)
)

const (
	LoggerType = "default"
)

func init() {
	RegisterLogger(LoggerType, DefaultLog)
}

//RegisterLogger 暴露函数覆盖默认logger
func RegisterLogger(name string, log logger.ILogger) {
	if _logMap == nil {
		_logMap = make(map[string]logger.ILogger)
	}
	_logMap[name] = log
	log.Init()
}

func GetLogger(name string) logger.ILogger {
	if l, ok := _logMap[name]; ok {
		return l
	}
	return DefaultLog
}
