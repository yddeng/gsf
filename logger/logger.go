package logger

var logger *Logger

func New(basePath string, fileName string) *Logger {
	logger_ := NewLogger(basePath, fileName, 1024*1024*2)
	logger_.Debugf("%s logger init", fileName)
	return logger_
}

func InitLogger(logger_ *Logger) {
	logger = logger_
}
