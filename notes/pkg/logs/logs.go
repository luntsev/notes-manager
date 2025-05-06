package logs

import (
	"log"
	"notes-manager/configs"
	"notes-manager/pkg/enum"
)

type Logger struct {
	logLevel int
}

func NewLogger(conf *configs.Config) *Logger {
	return &Logger{
		logLevel: conf.LogLevel,
	}
}

func (l *Logger) WriteToLog(msg string, msgType int) {
	switch l.logLevel {
	case enum.Debug:
		log.Println(msg)
	case enum.Normal:
		if msgType == enum.ErrorMsg || msgType == enum.WarningMsg {
			log.Println(msg)
		}
	case enum.Silent:
		if msgType == enum.ErrorMsg {
			log.Println(msg)
		}
	default:
		log.Println("logging level is undefined:", msg)
	}
}
