package routes

import (
	"book-manager/internal/handlers"

	"github.com/labstack/echo/v4"
)

type CategoryRoutes struct {
	Handler *handlers.CategoryHandler
}

func NewCategoryRoutes(h *handlers.CategoryHandler) *CategoryRoutes {
	return &CategoryRoutes{Handler: h}
}

func (r *CategoryRoutes) register(group *echo.Group) {
	c := group.Group("/v1/categories")

	c.POST("/add", r.Handler.Create)
	c.POST("/get-one", r.Handler.GetOne)
}
