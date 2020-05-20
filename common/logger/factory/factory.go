/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package factory

import (
	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/liangjfblue/gpusher/common/logger/vlog"
)

type IFactory interface {
	CreateLog(...logger.Option) logger.ILogger
}

//VLogFactor vlog factory
type VLogFactor struct{}

func (l *VLogFactor) CreateLog(opts ...logger.Option) logger.ILogger {
	return vlog.New(opts...)
}
