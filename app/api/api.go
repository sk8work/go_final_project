package api

import (
	"net/http"
)

// Init регистрирует все API обработчики
func Init() {
	// Регистрируем обработчики
	http.HandleFunc("/api/nextdate", NextDateHandler)

	// Здесь будут регистрироваться другие обработчики на следующих шагах
}
