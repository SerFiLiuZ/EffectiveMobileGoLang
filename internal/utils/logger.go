package utils

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	Debug bool
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC),
		Debug:  false,
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf("[INFO]  "+format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Fatalf("[ERROR] "+format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.Debug {
		l.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) EnableDebug() {
	l.Debug = true
}

func (l *Logger) DisableDebug() {
	l.Debug = false
}
