package routes

import (
	"book-manager/internal/handlers"
	"book-manager/internal/middleware"
	"book-manager/internal/repositories"
	"book-manager/internal/services/impl"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func registerAuthRoutes(api *echo.Group, db *gorm.DB) {
	authRepo := &repositories.AuthRepository{DB: db}
	authService := impl.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)
	authRoute := NewAuthRoutes(authHandler)
	authRoute.register(api)
}

func registerCategoryRoutes(api *echo.Group, db *gorm.DB) {
	categoryRepo := &repositories.CategoryRepository{DB: db}
	categoryService := impl.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	categoryRoute := NewCategoryRoutes(categoryHandler)
	categoryRoute.register(api)
}

func registerAuthorRoutes(api *echo.Group, db *gorm.DB) {
	authorRepo := &repositories.AuthorRepository{DB: db}
	authorService := impl.NewAuthorService(authorRepo)
	authorHandler := handlers.NewAuthorHandler(authorService)
	authorRoute := NewAuthorRoutes(authorHandler)
	authorRoute.register(api)
}

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/book-store/api")

	// Public routes (no JWT required)
	registerAuthRoutes(api, db)

	// Protected routes (JWT required)
	protected := api.Group("", middleware.JWTAuth())
	registerCategoryRoutes(protected, db)
	registerAuthorRoutes(protected, db)
}
