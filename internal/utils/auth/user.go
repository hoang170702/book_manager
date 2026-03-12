package auth

import "github.com/labstack/echo/v4"

// GetCurrentUser extracts the currently authenticated user from the Echo context.
// Currently returns "Anonymous" as a placeholder until JWT authentication is implemented.
func GetCurrentUser(c echo.Context) string {
	// TODO: Extract user from JWT token claims (e.g., c.Get("user"))
	return "Anonymous"
}
