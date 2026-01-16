package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sk8work/go_final_project/app/config"
)

// JWTSecretKey ключ для подписи JWT
var JWTSecretKey = []byte("todo-secret-key")
var cfg *config.Config

// Init инициализирует модуль аутентификации
func Init(config *config.Config) {
	cfg = config

	if cfg.IsAuthEnabled() {
		fmt.Printf("Аутентификация включена. Хэш пароля: %s\n", cfg.PasswordHash[:16]+"...")
	} else {
		fmt.Println("Аутентификация отключена (TODO_PASSWORD не установлен)")
	}
}

// GenerateToken создает JWT токен
func GenerateToken() (string, error) {
	if cfg == nil || !cfg.IsAuthEnabled() {
		return "", fmt.Errorf("аутентификация не настроена")
	}

	// Создаем claims с хэшем пароля
	claims := jwt.MapClaims{
		"password_hash": cfg.PasswordHash,
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
		"iat":           time.Now().Unix(),
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании токена: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет JWT токен
func ValidateToken(tokenString string) (bool, error) {
	if cfg == nil || !cfg.IsAuthEnabled() {
		// Если аутентификация отключена, пропускаем
		return true, nil
	}

	if tokenString == "" {
		return false, fmt.Errorf("токен не предоставлен")
	}

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("ошибка при парсинге токена: %w", err)
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем хэш пароля
		if hash, ok := claims["password_hash"].(string); ok {
			return hash == cfg.PasswordHash, nil
		}
		return false, fmt.Errorf("неверный формат claims")
	}

	return false, fmt.Errorf("недействительный токен")
}

// GetTokenFromRequest извлекает токен из запроса
func GetTokenFromRequest(r *http.Request) string {
	// Пробуем получить токен из куки
	if cookie, err := r.Cookie("token"); err == nil {
		return cookie.Value
	}

	// Пробуем получить токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Формат: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return ""
}

// IsEnabled проверяет, включена ли аутентификация
func IsEnabled() bool {
	return cfg != nil && cfg.IsAuthEnabled()
}
