package main

import (
	"book-manager/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	database.Connect()
}
