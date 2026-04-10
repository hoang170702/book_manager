package routes

import (
	"book-manager/internal/handlers"
	"book-manager/internal/middleware"
	"book-manager/internal/repositories"
	"book-manager/internal/services"
	"book-manager/internal/services/impl"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type authComponents struct {
	service services.IAuthService
	handler *handlers.AuthHandler
}

func registerAuthRoutes(api *echo.Group, db *gorm.DB) authComponents {
	authRepo := &repositories.AuthRepository{DB: db}
	authService := impl.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := api.Group("/auth")

	// Rate limit: 10 requests per minute per IP for login/register
	authLimited := auth.Group("", middleware.RateLimit(10, 1*time.Minute))
	authLimited.POST("/register", authHandler.Register)
	authLimited.POST("/login", authHandler.Login)

	// Refresh is public (needs refresh token, not access token)
	auth.POST("/refresh", authHandler.RefreshToken)

	return authComponents{service: authService, handler: authHandler}
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
	auth := registerAuthRoutes(api, db)

	// Protected routes (JWT required + token blacklist check)
	protected := api.Group("", middleware.JWTAuth(auth.service.IsTokenRevoked))
	registerCategoryRoutes(protected, db)
	registerAuthorRoutes(protected, db)

	// Logout is protected — requires access token to identify which tokens to blacklist
	protected.POST("/auth/logout", auth.handler.Logout)
}
