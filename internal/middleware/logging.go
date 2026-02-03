package middleware

import (
	"book-manager/internal/utils/logger"
	"time"

	"github.com/labstack/echo/v4"
)

// Logger middleware logs request/response information
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			requestID := GetRequestID(c)

			// Log incoming request
			logger.Info("HTTP", nil, "[%s] --> %s %s", requestID, req.Method, req.URL.Path)

			// Process request
			err := next(c)

			// Calculate latency
			latency := time.Since(start)
			status := c.Response().Status

			// Log response
			logger.Info("HTTP", nil, "[%s] <-- %s %s | %d | %v",
				requestID, req.Method, req.URL.Path, status, latency)

			return err
		}
	}
}
