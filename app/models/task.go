package models

import (
	"fmt"
	"time"
)

const (
	DateFormat      = "20060102"
	MaxTitleLength  = 255
	MaxRepeatLength = 128
)

// Task представляет собой задачу в планировщике
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"` // Формат: YYYYMMDD (20060102)
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"` // Правила повторения
}

// Validate проверяет корректность данных задачи
func (t *Task) Validate() error {
	// Проверяем формат даты
	if t.Date != "" {
		if _, err := time.Parse(DateFormat, t.Date); err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
	}

	return nil
}
