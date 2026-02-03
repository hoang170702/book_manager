package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const RequestIDKey = "request_id"

// RequestID middleware generates or extracts request ID from header
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID is in header
			requestID := c.Request().Header.Get("X-Request-ID")

			// If not, generate a new one
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Store in context
			c.Set(RequestIDKey, requestID)

			// Set in response header for traceability
			c.Response().Header().Set("X-Request-ID", requestID)

			return next(c)
		}
	}
}

// GetRequestID retrieves the request ID from echo context
func GetRequestID(c echo.Context) string {
	if id, ok := c.Get(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
