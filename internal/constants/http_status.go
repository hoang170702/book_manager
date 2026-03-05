package constants

import "net/http"

// HTTP Status Codes
const (
	StatusOK                  = http.StatusOK                  // 200
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusInternalServerError = http.StatusInternalServerError // 500
)
