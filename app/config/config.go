package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port   int
	WebDir string
	DBFile string
}

func Load() *Config {
	cfg := &Config{
		Port:   7540, // порт по умолчанию
		WebDir: "./web",
		DBFile: "scheduler.db", // путь к БД по умолчанию
	}

	// Получаем порт из переменной окружения TODO_PORT
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Port = port
		}
	}

	// Получаем путь к БД из переменной окружения TODO_DBFILE
	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		cfg.DBFile = dbFile
	}

	return cfg
}

// GetDBPath возвращает абсолютный путь к файлу БД
func (c *Config) GetDBPath() (string, error) {
	if filepath.IsAbs(c.DBFile) {
		return c.DBFile, nil
	}
	return filepath.Abs(c.DBFile)
}
