// logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Logger interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
	Close() error
}

type FileLogger struct {
	logger *log.Logger
	file   *os.File
	mu     sync.Mutex
}

// NewFileLogger creates a new FileLogger instance
func NewFileLogger(logPath string) (*FileLogger, error) {
	// Ensure the logs directory exists
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open the log file with append mode
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create the logger with timestamp
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return &FileLogger{
		logger: logger,
		file:   file,
	}, nil
}

func (l *FileLogger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[ERROR] "+format, v...)
}

func (l *FileLogger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *FileLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// LogEmailError logs email-related errors
func (l *FileLogger) LogEmailError(err error) {
	l.Error("Email error: %v", err)
}

// LogEmailSuccess logs successful email sending
func (l *FileLogger) LogEmailSuccess(recipient string) {
	l.Info("Email sent successfully to %s", recipient)
}
