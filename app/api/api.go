package api

import (
	"net/http"
)

// TaskHandler обрабатывает все методы для /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		// Для других методов вернем ошибку (будут реализованы позже)
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
	}
}
