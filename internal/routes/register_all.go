package routes

import (
	"book-manager/internal/handlers"
	"book-manager/internal/repositories"
	"book-manager/internal/services/impl"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func registerCategoryRoutes(api *echo.Group, db *gorm.DB) {
	categoryRepo := &repositories.CategoryRepository{DB: db}
	categoryService := impl.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	categoryRoute := NewCategoryRoutes(categoryHandler)
	categoryRoute.register(api)
}

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/book-store/api")

	// Đăng ký các routes theo module
	registerCategoryRoutes(api, db)
	// registerAuthorRoutes(api, db) // TODO: Thêm khi có author routes
	// registerBookRoutes(api, db)   // TODO: Thêm khi có book routes
}
