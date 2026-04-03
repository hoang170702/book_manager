package middleware

import (
	"book-manager/internal/constants"
	"book-manager/internal/utils"
	"book-manager/internal/utils/auth"
	"book-manager/internal/utils/enums/error_codes"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTAuth returns a middleware that validates JWT access tokens
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				resp := utils.BuildResponse[any](nil, error_codes.Unauthorized, "")
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			// Expect "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				resp := utils.BuildResponse[any](nil, error_codes.Unauthorized, "")
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			tokenString := parts[1]

			// Validate token
			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				resp := utils.BuildResponse[any](nil, error_codes.InvalidToken, "")
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			// Only accept access tokens, not refresh tokens
			if claims.Type != "access" {
				resp := utils.BuildResponse[any](nil, error_codes.InvalidToken, "")
				return c.JSON(constants.StatusUnauthorized, resp)
			}

			// Set claims in context for downstream handlers
			c.Set("user", claims)

			return next(c)
		}
	}
}
