package auth

import "github.com/labstack/echo/v4"

// GetCurrentUser extracts the currently authenticated user from the Echo context.
// Returns the username from JWT claims, or "Anonymous" if no valid claims are found.
func GetCurrentUser(c echo.Context) string {
	claims, ok := c.Get("user").(*Claims)
	if !ok {
		return "Anonymous"
	}
	return claims.Username
}
