package main

import (
	"book-manager/internal/config"
	"book-manager/internal/handlers"
	"book-manager/internal/middleware"
	"book-manager/internal/routes"
	"book-manager/internal/utils/logger"
	"book-manager/pkg/database"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func setupMiddleware(e *echo.Echo) {
	// Order matters: CORS first (handle preflight), then Recovery, RequestID, Logger
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.Recovery())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
}

func setupRoutes(e *echo.Echo, db *gorm.DB) {
	// Health check (no auth required)
	e.GET("/health", handlers.HealthCheck(db))

	// API routes
	routes.RegisterRoutes(e, db)
}

func startServer(e *echo.Echo, port string, db *gorm.DB) {
	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%s", port)
		logger.Info("SERVER", nil, "Starting server on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("SERVER", nil, "Shutting down server...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
		logger.Info("SERVER", nil, "Database connection closed")
	}

	logger.Info("SERVER", nil, "Server gracefully stopped")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("CONFIG", nil, ".env file not found, using environment variables")
	}

	cfg := config.LoadConfig()
	db := database.Connect()

	// Run migrations
	config.RunMigrations(db)

	e := echo.New()
	e.HideBanner = true // Hide default Echo banner for cleaner logs

	middleware.RegisterValidator(e)
	setupMiddleware(e)
	setupRoutes(e, db)

	startServer(e, cfg.Port, db)
}
