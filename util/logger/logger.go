package logger

import "github.com/yddeng/dutil/log"

var (
	logger  *log.Logger
	Debugln func(v ...interface{})
	Debugf  func(format string, v ...interface{})
	Infoln  func(v ...interface{})
	Infof   func(format string, v ...interface{})
	Errorln func(v ...interface{})
	Errorf  func(format string, v ...interface{})
)

func New(basePath string, fileName string) *log.Logger {
	logger_ := log.NewLogger(basePath, fileName, 1024*1024*2)
	logger_.Debugf("%s logger init", fileName)
	return logger_
}

func InitLogger(logger_ *log.Logger) {
	logger = logger_
	Debugln = logger_.Debugln
	Debugf = logger_.Debugf
	Infof = logger_.Infof
	Infoln = logger_.Infoln
	Errorf = logger_.Errorf
	Errorln = logger_.Errorln
}
