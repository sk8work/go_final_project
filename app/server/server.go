package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sk8work/go_final_project/app/api"
	"github.com/sk8work/go_final_project/app/auth"
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
	// Инициализируем аутентификацию с конфигурацией
	auth.Init(s.cfg)

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
	log.Printf("Database file: %s", s.cfg.DBFile)
	log.Printf("Serving static files from: %s", s.cfg.WebDir)

	// Проверяем настройки аутентификации
	if s.cfg.IsAuthEnabled() {
		log.Printf("Аутентификация включена")
		log.Printf("Для входа откройте: http://localhost:%d/login.html", s.cfg.Port)
	} else {
		log.Printf("Аутентификация отключена")
	}

	log.Printf("API endpoints available:")
	log.Printf("  POST   /api/signin - вход (если включена аутентификация)")
	log.Printf("  GET    /api/task - получение задачи")
	log.Printf("  POST   /api/task - добавление задачи")
	log.Printf("  PUT    /api/task - обновление задачи")
	log.Printf("  DELETE /api/task - удаление задачи")
	log.Printf("  POST   /api/task/done - завершение задачи")
	log.Printf("  GET    /api/nextdate - вычисление следующей даты")
	log.Printf("  GET    /api/tasks - получение списка задач")

	return s.httpServer.ListenAndServe()
}

// registerAPIRoutes регистрирует все API маршруты
func (s *Server) registerAPIRoutes(r chi.Router) {
	// Публичные маршруты (не требуют аутентификации)
	r.Post("/api/signin", api.SigninHandler)
	r.Get("/api/nextdate", api.NextDateHandler)

	// Защищенные маршруты (требуют аутентификации, если она включена)
	r.Route("/api", func(r chi.Router) {
		// Применяем middleware аутентификации ко всем маршрутам
		r.Use(authMiddleware)

		r.Post("/task", api.AddTaskHandler)
		r.Get("/task", api.GetTaskHandler)
		r.Put("/task", api.UpdateTaskHandler)
		r.Delete("/task", api.DeleteTaskHandler)
		r.Post("/task/done", api.DoneTaskHandler)
		r.Get("/tasks", api.TasksHandler)
	})
}

// authMiddleware middleware для проверки аутентификации
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Используем предварительно загруженную конфигурацию
		if !auth.IsEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из запроса
		token := auth.GetTokenFromRequest(r)

		// Проверяем токен
		valid, err := auth.ValidateToken(token)
		if !valid || err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Требуется аутентификация",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
