package logger

import (
	"context"
	"fmt"
	"log"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

// contextKey is a private type to prevent collisions in context
type contextKey string

const requestIDCtxKey contextKey = "gorm_request_id"

// ContextWithRequestID creates a context carrying a request ID for GORM logging
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, requestID)
}

// RequestIDFromContext extracts the request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return id
	}
	return ""
}

// GormLogger is a custom GORM logger that includes request ID in log output
type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

// NewGormLogger creates a new GORM logger with request ID support
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormlogger.Warn,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		requestID := RequestIDFromContext(ctx)
		logGorm("INFO", requestID, msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		requestID := RequestIDFromContext(ctx)
		logGorm("WARN", requestID, msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		requestID := RequestIDFromContext(ctx)
		logGorm("ERROR", requestID, msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	requestID := RequestIDFromContext(ctx)

	switch {
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		logGorm("WARN", requestID, "[SLOW SQL >= %v] [%.3fms] [rows:%d] %s", l.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case l.LogLevel >= gormlogger.Info:
		logGorm("INFO", requestID, "[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

func logGorm(level string, requestID string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	if requestID != "" {
		log.Printf("[%s] [%s] [GORM] [RequestId: %s] %s", timestamp, level, requestID, message)
	} else {
		log.Printf("[%s] [%s] [GORM] %s", timestamp, level, message)
	}
}
