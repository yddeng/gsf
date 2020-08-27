package util

import (
	"github.com/yddeng/gsf/util/log"
)

var logger *log.Logger

func InitLogger(basePath string, fileName string, fileMax int) {
	logger = log.NewLogger(basePath, fileName, fileMax)
	logger.Debugf("%s logger init", fileName)
}

func Logger() *log.Logger {
	return logger
}
