package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sk8work/go_final_project/app/config"
	"github.com/sk8work/go_final_project/app/db"
	"github.com/sk8work/go_final_project/app/server"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Получаем абсолютный путь к файлу БД
	dbPath, err := cfg.GetDBPath()
	if err != nil {
		log.Fatalf("Failed to get database path: %v", err)
	}

	// Инициализация базы данных
	log.Printf("Initializing database at: %s", dbPath)
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Создание сервера
	srv := server.New(cfg)

	// Каналы для graceful shutdown
	done := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	log.Printf("Server started on port %d", cfg.Port)
	log.Printf("Database file: %s", dbPath)
	log.Printf("Serving static files from: %s", cfg.WebDir)

	// Информация об аутентификации
	if os.Getenv("TODO_PASSWORD") != "" {
		log.Printf("Аутентификация включена. Пароль установлен.")
		log.Printf("Для входа откройте: http://localhost:%d/login.html", cfg.Port)
	} else {
		log.Printf("Аутентификация отключена (TODO_PASSWORD не установлен)")
	}

	// Ожидание сигнала завершения или ошибки
	select {
	case <-done:
		log.Println("Server stopping...")
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	// Закрываем базу данных
	db.Close()
	log.Println("Server stopped")
}
