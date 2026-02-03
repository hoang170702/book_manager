package handlers

import (
	"book-manager/pkg/database"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// HealthCheck returns the health status of the application
func HealthCheck(c echo.Context) error {
	dbStatus := "connected"

	// Check database connection
	sqlDB, err := database.DB.DB()
	if err != nil {
		dbStatus = "disconnected"
	} else if err = sqlDB.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	status := "healthy"
	if dbStatus == "disconnected" {
		status = "unhealthy"
	}

	return c.JSON(200, HealthResponse{
		Status:   status,
		Database: dbStatus,
	})
}
