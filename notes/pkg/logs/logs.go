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

func (l *Logger) WriteInfo(msg string) {
	if l.logLevel == enum.Debug {
		log.Println("INFO:", msg)
	}
}

func (l *Logger) WriteWarn(msg string) {
	if l.logLevel == enum.Debug || l.logLevel == enum.Normal {
		log.Println("WARNING:", msg)
	}
}

func (l *Logger) WriteError(msg string) {
	log.Println("ERROR:", msg)
}
