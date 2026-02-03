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
)

func setupMiddleware(e *echo.Echo) {
	// Order matters: Recovery first, then RequestID, then Logger
	e.Use(middleware.Recovery())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
}

func setupRoutes(e *echo.Echo) {
	// Health check (no auth required)
	e.GET("/health", handlers.HealthCheck)

	// API routes
	routes.RegisterRoutes(e, database.DB)
}

func startServer(e *echo.Echo, port string) {
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
	if sqlDB, err := database.DB.DB(); err == nil {
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
	database.Connect()

	e := echo.New()
	e.HideBanner = true // Hide default Echo banner for cleaner logs

	setupMiddleware(e)
	setupRoutes(e)

	startServer(e, cfg.Port)
}
