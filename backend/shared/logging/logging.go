package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger interface for structured logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// StructuredLogger implements the Logger interface using logrus
type StructuredLogger struct {
	logger *logrus.Logger
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(service string) *StructuredLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Add service name to all log entries
	logger.AddHook(&ServiceHook{Service: service})

	return &StructuredLogger{
		logger: logger,
	}
}

// ServiceHook adds service name to log entries
type ServiceHook struct {
	Service string
}

func (h *ServiceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ServiceHook) Fire(entry *logrus.Entry) error {
	entry.Data["service"] = h.Service
	return nil
}

// Info logs an info message
func (l *StructuredLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.WithFields(parseArgs(args...)).Info(msg)
	} else {
		l.logger.Info(msg)
	}
}

// Error logs an error message
func (l *StructuredLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.WithFields(parseArgs(args...)).Error(msg)
	} else {
		l.logger.Error(msg)
	}
}

// Warn logs a warning message
func (l *StructuredLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.WithFields(parseArgs(args...)).Warn(msg)
	} else {
		l.logger.Warn(msg)
	}
}

// Debug logs a debug message
func (l *StructuredLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.WithFields(parseArgs(args...)).Debug(msg)
	} else {
		l.logger.Debug(msg)
	}
}

// Fatal logs a fatal message and exits
func (l *StructuredLogger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.WithFields(parseArgs(args...)).Fatal(msg)
	} else {
		l.logger.Fatal(msg)
	}
}

// parseArgs converts variadic arguments to logrus fields
func parseArgs(args ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			if key, ok := args[i].(string); ok {
				fields[key] = args[i+1]
			}
		}
	}
	
	return fields
}
