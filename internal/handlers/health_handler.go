package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// HealthCheck returns the health status of the application
func HealthCheck(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		dbStatus := "connected"

		// Check database connection
		sqlDB, err := db.DB()
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
}
