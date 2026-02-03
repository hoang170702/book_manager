package logger

import (
	"fmt"
	"log"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
)

// RequestIdGetter interface for extracting requestId
type RequestIdGetter interface {
	GetRequestId() string
}

// Info logs an informational message with optional request
func Info(module string, request RequestIdGetter, format string, args ...interface{}) {
	requestId := extractRequestId(request)
	logMessage(INFO, module, requestId, format, args...)
}

// Error logs an error message with optional request
func Error(module string, request RequestIdGetter, format string, args ...interface{}) {
	requestId := extractRequestId(request)
	logMessage(ERROR, module, requestId, format, args...)
}

// Warn logs a warning message with optional request
func Warn(module string, request RequestIdGetter, format string, args ...interface{}) {
	requestId := extractRequestId(request)
	logMessage(WARN, module, requestId, format, args...)
}

// Debug logs a debug message with optional request
func Debug(module string, request RequestIdGetter, format string, args ...interface{}) {
	requestId := extractRequestId(request)
	logMessage(DEBUG, module, requestId, format, args...)
}

// extractRequestId safely extracts requestId from request
func extractRequestId(request RequestIdGetter) string {
	if request == nil {
		return ""
	}
	return request.GetRequestId()
}

// logMessage is the internal function that formats and logs messages
func logMessage(level LogLevel, module string, requestId string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	if requestId != "" {
		log.Printf("[%s] [%s] [%s] [RequestId: %s] %s", timestamp, level, module, requestId, message)
	} else {
		log.Printf("[%s] [%s] [%s] %s", timestamp, level, module, message)
	}
}
