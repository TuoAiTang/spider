package log

import (
	l "log"
)

func Info(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func Error(format string, args ...interface{}) {
	l.Printf(format, args...)
}
