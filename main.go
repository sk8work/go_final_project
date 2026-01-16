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

	// Каналы для управления
	done := make(chan os.Signal, 1)
	serverError := make(chan error, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Starting server on port %d...", cfg.Port)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			// Отправляем ошибку в канал вместо log.Fatal
			serverError <- err
		} else {
			// Если сервер корректно остановлен
			serverError <- nil
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

	// Ожидание сигнала завершения или ошибки сервера
	select {
	case err := <-serverError:
		// Ошибка при запуске/работе сервера
		log.Printf("Server error: %v", err)
		// Отправляем сигнал для graceful shutdown
		done <- syscall.SIGTERM

	case <-done:
		// Получен сигнал завершения (Ctrl+C, SIGTERM)
		log.Println("Received shutdown signal")
	}

	log.Println("Server stopping...")

	// Graceful shutdown - сначала сервер, потом БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	// Закрываем соединение с БД после остановки сервера
	log.Println("Closing database connection...")
	if err := db.Close(); err != nil {
		log.Printf("Failed to close database: %v", err)
	}

	log.Println("Server stopped")
}
