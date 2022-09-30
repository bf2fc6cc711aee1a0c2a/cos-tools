package log

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/redhat-developer/app-services-cli/pkg/core/logging"
)

const debugLevel = 1
const infoLevel = 0

type GLogLogger struct {
}

func (l *GLogLogger) DebugEnabled() bool {
	return bool(glog.V(debugLevel))
}

func (l *GLogLogger) InfoEnabled() bool {
	return bool(glog.V(infoLevel))
}

func (l *GLogLogger) SetDebug(_ bool) {
	panic("not implemented")
}

func (l *GLogLogger) Debug(args ...interface{}) {
	if glog.V(debugLevel) {
		glog.InfoDepth(1, args...)
	}

}
func (l *GLogLogger) Debugf(format string, args ...interface{}) {
	if glog.V(debugLevel) {
		msg := fmt.Sprintf(format, args...)
		glog.InfoDepth(1, msg)
	}
}

func (l *GLogLogger) Info(args ...interface{}) {
	if glog.V(infoLevel) {
		glog.InfoDepth(1, args...)
	}

}
func (l *GLogLogger) Infof(format string, args ...interface{}) {
	if glog.V(infoLevel) {
		msg := fmt.Sprintf(format, args...)
		glog.InfoDepth(1, msg)
	}
}

func (l *GLogLogger) Error(args ...interface{}) {
	if glog.V(infoLevel) {
		glog.ErrorDepth(1, args...)
	}
}

func (l *GLogLogger) Errorf(format string, args ...interface{}) {
	if glog.V(infoLevel) {
		msg := fmt.Sprintf(format, args...)
		glog.ErrorDepth(1, msg)
	}
}

func NewLogger() logging.Logger {
	return &GLogLogger{}
}
