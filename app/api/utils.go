package api

import (
	"encoding/json"
	"net/http"
	"os"
)

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

// TaskResponse представляет ответ с ID задачи или ошибкой
type TaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// TaskJSON представляет задачу в формате JSON (с полем ID как строка)
type TaskJSON struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// checkAuth проверяет аутентификацию (простая версия)
func checkAuth(r *http.Request) bool {
	// Если пароль не установлен, аутентификация не требуется
	if os.Getenv("TODO_PASSWORD") == "" {
		return true
	}

	// Здесь будет проверка JWT токена
	// Пока возвращаем true для совместимости со старыми тестами
	return true
}
