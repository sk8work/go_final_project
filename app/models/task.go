package models

import (
	"fmt"
	"time"
)

// Task представляет собой задачу в планировщике
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"` // Формат: YYYYMMDD (20060102)
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"` // Правила повторения (до 128 символов)
}

// Validate проверяет корректность данных задачи
func (t *Task) Validate() error {
	// Проверяем формат даты
	if t.Date != "" {
		if _, err := time.Parse("20060102", t.Date); err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
	}

	// Проверяем длину заголовка
	if len(t.Title) > 255 {
		return fmt.Errorf("title too long, max 255 characters")
	}

	// Проверяем длину правила повторения
	if len(t.Repeat) > 128 {
		return fmt.Errorf("repeat rule too long, max 128 characters")
	}

	return nil
}

// ToMap преобразует задачу в map для SQL запросов
func (t *Task) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"date":    t.Date,
		"title":   t.Title,
		"comment": t.Comment,
		"repeat":  t.Repeat,
	}
}
