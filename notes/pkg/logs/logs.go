package logs

import (
	"github.com/luntsev/notes-manager/notes/pkg/enum"
	"log"
)

type Logger struct {
	logLevel int
}

func NewLogger(level int) *Logger {
	return &Logger{
		logLevel: level,
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
