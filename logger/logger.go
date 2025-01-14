// logger/logger.go
package logger

import "log"

type Logger interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
}

type FileLogger struct {
	logger *log.Logger
}
