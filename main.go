package main

import (
	"context"
	"github.com/sk8work/go_final_project/app/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sk8work/go_final_project/app/config"
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
	defer db.Close()

	// Создание сервера
	srv := server.New(cfg)

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %d", cfg.Port)
	log.Printf("Database file: %s", dbPath)

	// Ожидание сигнала завершения
	<-done
	log.Println("Server stopping...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
