package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port         int
	WebDir       string
	DBFile       string
	Password     string
	PasswordHash string
	AuthEnabled  bool
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

	// Получаем пароль из переменной окружения TODO_PASSWORD
	if password := os.Getenv("TODO_PASSWORD"); password != "" {
		cfg.Password = password
		cfg.AuthEnabled = true

		// Вычисляем хэш пароля один раз при загрузке конфигурации
		hash := sha256.Sum256([]byte(password))
		cfg.PasswordHash = hex.EncodeToString(hash[:])
	} else {
		cfg.AuthEnabled = false
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

// IsAuthEnabled проверяет, включена ли аутентификация
func (c *Config) IsAuthEnabled() bool {
	return c.AuthEnabled
}
