package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Регистрация маршрутов
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
	log.Printf("Serving static files from: %s", s.cfg.WebDir)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
