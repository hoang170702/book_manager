package main

import (
	"book-manager/internal/config"
	"book-manager/internal/routes"
	"book-manager/pkg/database"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func startServer(port string) {
	e := echo.New()

	routes.RegisterRoutes(e, database.DB)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	cfg := config.LoadConfig()
	database.Connect()

	startServer(cfg.Port)
}
