package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type CustomLogger struct {
	mu sync.Mutex
	*log.Logger
	level string
}

// NewCustomLogger creates a new CustomLogger
func NewLogger(out io.Writer, prefix string, flag int) *CustomLogger {
	return &CustomLogger{
		Logger: log.New(out, prefix, flag),
		level:  "INFO",
	}
}

func (l *CustomLogger) Debug(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.setPrefix("DEBUG")
	l.Println(v...)
}

func (l *CustomLogger) Info(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.setPrefix("INFO")
	l.Println(v...)
}

func (l *CustomLogger) Warn(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.setPrefix("WARN")
	l.Println(v...)
}

func (l *CustomLogger) Error(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.setPrefix("ERROR")
	l.Println(v...)
}

func (l *CustomLogger) setPrefix(level string) {
	l.SetPrefix(fmt.Sprintf("%s: ", level))
}

func setupLogger(logFile string) (*CustomLogger, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	return NewLogger(mw, "", log.LstdFlags), nil
}
