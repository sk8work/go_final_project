package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sk8work/go_final_project/app/api"
	"github.com/sk8work/go_final_project/app/config"
	"github.com/sk8work/go_final_project/app/handlers"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Регистрируем API маршруты
	s.registerAPIRoutes(r)

	// Регистрация статических файлов
	handler := handlers.NewHandler(s.cfg.WebDir)
	handler.RegisterRoutes(r)

	// Настройка сервера
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %d", s.cfg.Port)
	log.Printf("API endpoints available:")
	log.Printf("  POST /api/task - добавление задачи")
	log.Printf("  GET  /api/nextdate - вычисление следующей даты")

	return s.httpServer.ListenAndServe()
}

// registerAPIRoutes регистрирует все API маршруты
func (s *Server) registerAPIRoutes(r chi.Router) {
	// API маршруты
	r.Route("/api", func(r chi.Router) {
		r.Get("/nextdate", api.NextDateHandler)
		r.Post("/task", api.TaskHandler)
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
