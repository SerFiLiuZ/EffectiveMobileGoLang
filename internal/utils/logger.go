package utils

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.LUTC),
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}
