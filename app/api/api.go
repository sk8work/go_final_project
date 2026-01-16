package api

import (
	"encoding/json"
	"net/http"
)

const (
	TasksLimit = 50
	DateFormat = "20060102"
)

// TaskResponse представляет ответ с ID задачи или ошибкой
type TaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// writeJSON записывает JSON ответ
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// writeJSONError записывает JSON с ошибкой
func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
}
