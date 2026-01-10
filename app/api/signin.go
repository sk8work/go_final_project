package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sk8work/go_final_project/app/auth"
)

// SigninRequest запрос на аутентификацию
type SigninRequest struct {
	Password string `json:"password"`
}

// SigninResponse ответ на аутентификацию
type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

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
		writeJSON(w, SigninResponse{Error: "ошибка разбора JSON: " + err.Error()})
		return
	}

	// Получаем пароль из переменной окружения
	envPassword := os.Getenv("TODO_PASSWORD")

	// Если пароль не установлен в окружении, возвращаем ошибку
	if envPassword == "" {
		writeJSON(w, SigninResponse{Error: "аутентификация не настроена"})
		return
	}

	// Проверяем пароль
	if req.Password != envPassword {
		writeJSON(w, SigninResponse{Error: "неверный пароль"})
		return
	}

	// Генерируем токен
	token, err := auth.GenerateToken()
	if err != nil {
		writeJSON(w, SigninResponse{Error: "ошибка при создании токена: " + err.Error()})
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
	writeJSON(w, SigninResponse{Token: token})
}
