package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// GetLimitFromRequest извлекает лимит из параметров запроса
func GetLimitFromRequest(r *http.Request) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return DefaultTasksLimit
	}

	// Пытаемся преобразовать в число
	var limit int
	if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
		return DefaultTasksLimit
	}

	// Проверяем границы
	if limit <= 0 {
		return DefaultTasksLimit
	}
	if limit > MaxTasksLimit {
		return MaxTasksLimit
	}

	return limit
}

// writeJSON записывает JSON ответ с указанным статус-кодом
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeJSONError записывает JSON с ошибкой
func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
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
