package logging
package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents a log level.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

// Logger provides logging functionality.
type Logger struct {
	level    Level
	file     *os.File
	handlers []Handler
}

// Handler is a custom log handler.
type Handler func(entry LogEntry)

// LogEntry represents a single log entry.
type LogEntry struct {
	Level     Level
	Message   string
	Context   map[string]interface{}
	Timestamp time.Time
}

// New creates a new logger instance.
func New(level Level) *Logger {
	return &Logger{
		level:    level,
		handlers: make([]Handler, 0),
	}
}

// SetFile sets the log file path.
func (l *Logger) SetFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.file = file
	return nil
}

// AddHandler adds a custom log handler.
func (l *Logger) AddHandler(h Handler) {
	l.handlers = append(l.handlers, h)
}

// log logs an entry at the given level.
func (l *Logger) log(level Level, message string, context map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Level:     level,
		Message:   message,
		Context:   context,
		Timestamp: time.Now(),
	}

	// Call handlers
	for _, handler := range l.handlers {
		handler(entry)
	}

	// Write to file if set
	if l.file != nil {
		formatted := l.formatEntry(entry)
		l.file.WriteString(formatted + "\n")
	}

	// Write to stdout in debug mode
	if level <= INFO {
		fmt.Println(l.formatEntry(entry))
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(DEBUG, message, ctx)
}

// Info logs an info message.
func (l *Logger) Info(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(INFO, message, ctx)
}

// Warning logs a warning message.
func (l *Logger) Warning(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(WARNING, message, ctx)
}

// Error logs an error message.
func (l *Logger) Error(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(ERROR, message, ctx)
}

// Critical logs a critical message.
func (l *Logger) Critical(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(CRITICAL, message, ctx)
}

// formatEntry formats a log entry for display.
func (l *Logger) formatEntry(entry LogEntry) string {
	levelStr := l.levelString(entry.Level)
	return fmt.Sprintf("[%s] %s: %s", entry.Timestamp.Format("2006-01-02 15:04:05"), levelStr, entry.Message)
}

// levelString returns the string representation of a level.
func (l *Logger) levelString(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Close closes the log file.
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// StdoutHandler writes log entries to stdout.
func StdoutHandler(entry LogEntry) {
	log.Printf("[%s] %s", entry.Timestamp.Format("15:04:05"), entry.Message)
}
