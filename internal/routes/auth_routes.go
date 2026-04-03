package routes

import (
	"book-manager/internal/handlers"

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
	auth.POST("/register", r.handler.Register)
	auth.POST("/login", r.handler.Login)
	auth.POST("/refresh", r.handler.RefreshToken)
}
