/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package vlog

import (
	"os"

	"github.com/liangjfblue/gpusher/common/logger"
	"github.com/sirupsen/logrus"
)

type VLog struct {
	opts logger.Options
}

func (l *VLog) Init(opts ...logger.Option) {
	for _, o := range opts {
		o(&l.opts)
	}

	l.initVLog()
}

func (l *VLog) initVLog() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.Level(l.opts.Level))
	logrus.SetReportCaller(false)
}

func (l *VLog) Trace(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

func (l *VLog) Debug(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func (l *VLog) Info(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func (l *VLog) Warn(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func (l *VLog) Error(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func New(opts ...logger.Option) *VLog {
	l := new(VLog)
	l.opts = logger.DefaultOptions

	for _, o := range opts {
		o(&l.opts)
	}

	return l
}
