package middleware

import (
	"book-manager/internal/constants"
	"book-manager/internal/utils"
	"book-manager/internal/utils/auth"
	"book-manager/internal/utils/enums/error_codes"
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

// TokenRevokeChecker is a function that checks if a token is revoked
type TokenRevokeChecker func(token string) bool

// extractRequestId reads request_id from JSON body without consuming it
func extractRequestId(c echo.Context) string {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return ""
	}
	// Put body back so handler can read it
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	var partial struct {
		RequestId string `json:"request_id"`
	}
	if err := json.Unmarshal(body, &partial); err == nil {
		return partial.RequestId
	}
	return ""
}

// JWTAuth returns a middleware that validates JWT access tokens.
// Optional: pass a TokenRevokeChecker to also check token blacklist.
func JWTAuth(checkers ...TokenRevokeChecker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try to extract request_id from body for error responses
			requestId := extractRequestId(c)

			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				resp := utils.BuildResponse[any](nil, error_codes.Unauthorized, requestId)
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			// Expect "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				resp := utils.BuildResponse[any](nil, error_codes.Unauthorized, requestId)
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			tokenString := parts[1]

			// Validate token
			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				resp := utils.BuildResponse[any](nil, error_codes.InvalidToken, requestId)
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			// Only accept access tokens, not refresh tokens
			if claims.Type != "access" {
				resp := utils.BuildResponse[any](nil, error_codes.InvalidToken, requestId)
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			// Check if token is blacklisted (revoked via logout)
			if len(checkers) > 0 && checkers[0] != nil {
				if checkers[0](tokenString) {
					resp := utils.BuildResponse[any](nil, error_codes.InvalidToken, requestId)
					return c.JSON(constants.StatusUnauthorized, resp)
				}
			}

			// Set claims in context for downstream handlers
			c.Set("user", claims)

			return next(c)
		}
	}
}
