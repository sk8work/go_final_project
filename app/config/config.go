package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port   int
	WebDir string
}

func Load() *Config {
	cfg := &Config{
		Port:   7540, // порт по умолчанию
		WebDir: "./web",
	}

	// Получаем порт из переменной окружения TODO_PORT
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	return cfg
}
