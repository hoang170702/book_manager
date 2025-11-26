package routes

import (
	"book-manager/internal/handlers"

	"github.com/labstack/echo/v4"
)

type AuthorRoutes struct {
	h *handlers.AuthorHandler
}

func NewAuthorRoutes(h *handlers.AuthorHandler) *AuthorRoutes {
	return &AuthorRoutes{h: h}
}

func (r *AuthorRoutes) register(group *echo.Group) {
	c := group.Group("/v1/authors")

	c.POST("/add", r.h.Add)
}
