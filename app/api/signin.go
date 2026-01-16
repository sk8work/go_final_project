package api

import (
	"encoding/json"
	"net/http"

	"github.com/sk8work/go_final_project/app/auth"
)

// SigninHandler обрабатывает запрос на вход
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON запрос
	var req SigninRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, SigninResponse{Error: "неверный формат JSON: " + err.Error()})
		return
	}

	// Используем auth.IsEnabled() вместо прямого обращения к os.Getenv
	if !auth.IsEnabled() {
		writeJSON(w, http.StatusServiceUnavailable, SigninResponse{Error: "аутентификация не настроена"})
		return
	}

	// Проверяем пароль через auth модуль (который использует config)
	token, err := auth.GenerateToken()
	if err != nil {
		// Если ошибка генерации токена, проверяем пароль
		writeJSON(w, http.StatusUnauthorized, SigninResponse{Error: "неверный пароль"})
		return
	}

	// Устанавливаем куку с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   8 * 60 * 60, // 8 часов
		HttpOnly: true,
		Secure:   false, // Для разработки, в production должно быть true
		SameSite: http.SameSiteStrictMode,
	})

	// Возвращаем токен
	writeJSON(w, http.StatusOK, SigninResponse{Token: token})
}
