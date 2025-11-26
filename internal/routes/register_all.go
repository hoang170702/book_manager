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

func registerAuthorRoutes(api *echo.Group, db *gorm.DB) {
	authorRepo := &repositories.AuthorRepository{DB: db}
	categoryService := impl.NewAuthorService(authorRepo)
	authorHandler := handlers.NewAuthorHandler(categoryService)
	authorRoute := NewAuthorRoutes(authorHandler)
	authorRoute.register(api)
}

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/book-store/api")

	registerCategoryRoutes(api, db)
	registerAuthorRoutes(api, db)
	// registerBookRoutes(api, db)   // TODO: Thêm khi có book routes
}
