package api

import (
	"fmt"
	"time"

	"github.com/sk8work/go_final_project/app/models"
)

// TaskRequest представляет запрос для создания или обновления задачи
type TaskRequest struct {
	ID      string `json:"id,omitempty"`      // Только для обновления
	Date    string `json:"date,omitempty"`    // Формат: 20060102
	Title   string `json:"title"`             // Обязательное поле
	Comment string `json:"comment,omitempty"` // Необязательное поле
	Repeat  string `json:"repeat,omitempty"`  // Правило повторения
}

// Validate проверяет корректность TaskRequest
func (r *TaskRequest) Validate() error {
	// Проверка формата даты
	if r.Date != "" {
		_, err := time.Parse(models.DateFormat, r.Date)
		if err != nil {
			return fmt.Errorf("неверный формат даты, ожидается %s", models.DateFormat)
		}
	}

	// Проверка обязательного поля
	if r.Title == "" {
		return fmt.Errorf("заголовок обязателен")
	}

	return nil
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

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Tasks []TaskJSON `json:"tasks"`
	Error string     `json:"error,omitempty"`
}

// SigninRequest запрос на аутентификацию
type SigninRequest struct {
	Password string `json:"password"`
}

// SigninResponse ответ на аутентификацию
type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}
