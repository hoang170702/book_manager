package middleware

import (
	"book-manager/internal/utils/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// bodyRequest extracts request_id from JSON body
type bodyRequest struct {
	RequestId string `json:"request_id"`
}

// bodyWriter wraps http.ResponseWriter to capture response body
type bodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger middleware logs request/response information with body content
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			headerRequestID := GetRequestID(c)

			// Read request body and restore it
			var reqBody []byte
			if req.Body != nil {
				reqBody, _ = io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			// Extract request_id from body if present
			bodyReqID := ""
			if len(reqBody) > 0 {
				var br bodyRequest
				if err := json.Unmarshal(reqBody, &br); err == nil && br.RequestId != "" {
					bodyReqID = br.RequestId
				}
			}

			// Use body request_id if available, fallback to header UUID
			logReqID := headerRequestID
			if bodyReqID != "" {
				logReqID = bodyReqID
			}

			// Log incoming request
			logger.Info("HTTP", nil, "[%s] --> %s %s | body: %s",
				logReqID, req.Method, req.URL.Path, string(reqBody))

			// Capture response body
			resBody := new(bytes.Buffer)
			writer := &bodyWriter{ResponseWriter: c.Response().Writer, body: resBody}
			c.Response().Writer = writer

			// Process request
			err := next(c)

			// Calculate latency
			latency := time.Since(start)
			status := c.Response().Status

			// Log response
			logger.Info("HTTP", nil, "[%s] <-- %s %s | %d | %v | body: %s",
				logReqID, req.Method, req.URL.Path, status, latency, resBody.String())

			return err
		}
	}
}
