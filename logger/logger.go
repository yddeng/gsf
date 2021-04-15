package logger

import "github.com/yddeng/dutil/log"

var (
	logger Interface
	Debug  func(v ...interface{})
	Debugf func(format string, v ...interface{})
	Info   func(v ...interface{})
	Infof  func(format string, v ...interface{})
	Error  func(v ...interface{})
	Errorf func(format string, v ...interface{})
)

type Interface interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

func New(basePath string, fileName string) *log.Logger {
	logger_ := log.NewLogger(basePath, fileName, 1024*1024*2)
	logger_.Debugf("%s logger init", fileName)
	return logger_
}

func InitLogger(logger_ Interface) {
	logger = logger_
	Debug = logger_.Debug
	Debugf = logger_.Debugf
	Infof = logger_.Infof
	Info = logger_.Info
	Errorf = logger_.Errorf
	Error = logger_.Error
}
