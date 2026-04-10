package config

import (
	"book-manager/internal/utils"
)

type Config struct {
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		Port:      utils.GetEnv("PORT", "8080"),
		JWTSecret: utils.GetEnv("JWT_SECRET", "book-manager-default-secret-change-me"),
	}
}
