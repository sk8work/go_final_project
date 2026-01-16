package models

import (
	"fmt"
	"time"
)

// Константы для форматов даты
const (
	DateFormat    = "20060102"   // Основной формат даты для базы данных
	DisplayFormat = "02.01.2006" // Формат для отображения пользователю
	ISOFormat     = "2006-01-02" // ISO формат
	TimeFormat    = "15:04"      // Формат времени
)

// Task представляет собой задачу в планировщике
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"` // Формат: YYYYMMDD (20060102)
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"` // Правила повторения (до 128 символов)
}

// Validate проверяет корректность данных задачи (базовая бизнес-логика)
func (t *Task) Validate() error {
	// Проверяем формат даты (это бизнес-логика, а не валидация БД)
	if t.Date != "" {
		if _, err := time.Parse(DateFormat, t.Date); err != nil {
			return fmt.Errorf("invalid date format, expected %s: %w", DateFormat, err)
		}
	}

	// Проверяем обязательные поля (бизнес-логика)
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}

	// Проверяем правило повторения (бизнес-логика)
	if t.Repeat != "" {
		if err := validateRepeatRule(t.Repeat); err != nil {
			return fmt.Errorf("invalid repeat rule: %w", err)
		}
	}

	return nil
}

// validateRepeatRule проверяет корректность правила повторения
func validateRepeatRule(repeat string) error {
	// Базовая проверка - не пустая строка и не слишком длинная
	// Детальная проверка выполняется в функции NextDate
	if repeat == "" {
		return nil // Пустое правило - это нормально
	}

	if len(repeat) > 128 {
		return fmt.Errorf("repeat rule too long, max 128 characters")
	}

	return nil
}

// ParseDate парсит строку даты в формате DateFormat
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// FormatDate форматирует время в строку DateFormat
func FormatDate(date time.Time) string {
	return date.Format(DateFormat)
}

// FormatDisplay форматирует дату для отображения пользователю
func FormatDisplay(dateStr string) (string, error) {
	date, err := ParseDate(dateStr)
	if err != nil {
		return "", err
	}
	return date.Format(DisplayFormat), nil
}

// ParseDisplay парсит дату из пользовательского формата
func ParseDisplay(displayStr string) (string, error) {
	date, err := time.Parse(DisplayFormat, displayStr)
	if err != nil {
		return "", err
	}
	return date.Format(DateFormat), nil
}

// IsValidDate проверяет, является ли строка корректной датой в основном формате
func IsValidDate(dateStr string) bool {
	_, err := ParseDate(dateStr)
	return err == nil
}

// Today возвращает текущую дату в основном формате
func Today() string {
	return FormatDate(time.Now())
}

// IsToday проверяет, является ли дата сегодняшней
func IsToday(dateStr string) bool {
	today := Today()
	return dateStr == today
}

// IsPastDate проверяет, находится ли дата в прошлом
func IsPastDate(dateStr string) bool {
	date, err := ParseDate(dateStr)
	if err != nil {
		return false
	}
	return date.Before(time.Now().Truncate(24 * time.Hour))
}
