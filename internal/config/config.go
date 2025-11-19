package config

import (
	"book-manager/internal/utils"
)

type Config struct {
	Port string
}

func LoadConfig() *Config {
	return &Config{
		Port: utils.GetEnv("PORT", "8080"),
	}
}
