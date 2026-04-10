package routes

import (
	"book-manager/internal/handlers"
	"book-manager/internal/middleware"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthRoutes struct {
	handler *handlers.AuthHandler
}

func NewAuthRoutes(handler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{handler: handler}
}

func (r *AuthRoutes) register(api *echo.Group) {
	auth := api.Group("/auth")

	// Rate limit: 10 requests per minute per IP for login/register
	authLimited := auth.Group("", middleware.RateLimit(10, 1*time.Minute))
	authLimited.POST("/register", r.handler.Register)
	authLimited.POST("/login", r.handler.Login)

	// Refresh doesn't need strict rate limiting
	auth.POST("/refresh", r.handler.RefreshToken)
}
