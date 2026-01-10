package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecretKey ключ для подписи JWT
var JWTSecretKey = []byte("todo-secret-key")

// PasswordHash хэш пароля из переменной окружения
var PasswordHash string

// Init инициализирует модуль аутентификации
func Init() {
	password := os.Getenv("TODO_PASSWORD")
	if password != "" {
		// Вычисляем хэш пароля
		hash := sha256.Sum256([]byte(password))
		PasswordHash = hex.EncodeToString(hash[:])
		fmt.Printf("Аутентификация включена. Хэш пароля: %s\n", PasswordHash[:16]+"...")
	}
}

// GenerateToken создает JWT токен
func GenerateToken() (string, error) {
	if PasswordHash == "" {
		return "", fmt.Errorf("аутентификация не настроена")
	}

	// Создаем claims с хэшем пароля
	claims := jwt.MapClaims{
		"password_hash": PasswordHash,
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
	if PasswordHash == "" {
		// Если пароль не установлен, аутентификация не требуется
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
			return hash == PasswordHash, nil
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

// Middleware создает middleware для проверки аутентификации
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если пароль не установлен, пропускаем без проверки
		if PasswordHash == "" {
			next(w, r)
			return
		}

		// Получаем токен из запроса
		token := GetTokenFromRequest(r)

		// Проверяем токен
		valid, err := ValidateToken(token)
		if !valid || err != nil {
			// Возвращаем ошибку аутентификации
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Требуется аутентификация",
			})
			return
		}

		// Токен валиден, продолжаем
		next(w, r)
	}
}
