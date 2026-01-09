package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка на пустое правило
	if repeat == "" {
		return "", errors.New("повторение не указано")
	}

	// Парсим начальную дату
	startDate, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %w", err)
	}

	// Разбиваем правило на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила повторения")
	}

	ruleType := parts[0]

	switch ruleType {
	case "d":
		// Правило: d <число>
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d: требуется число дней")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("неверное число дней")
		}

		if days <= 0 || days > 400 {
			return "", errors.New("интервал должен быть от 1 до 400 дней")
		}

		return calculateNextDateByDays(now, startDate, days), nil

	case "y":
		// Правило: y (ежегодно)
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y")
		}
		return calculateNextDateByYears(now, startDate), nil

	default:
		// Для остальных правил возвращаем ошибку (базовая реализация)
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// calculateNextDateByDays вычисляет следующую дату для правила d
func calculateNextDateByDays(now, startDate time.Time, interval int) string {
	date := startDate

	// Если начальная дата уже в будущем, возвращаем её
	if afterNow(date, now) {
		return date.Format(dateFormat)
	}

	// Вычисляем следующую дату, добавляя интервалы
	for !afterNow(date, now) {
		date = date.AddDate(0, 0, interval)
	}

	return date.Format(dateFormat)
}

// calculateNextDateByYears вычисляет следующую дату для правила y
func calculateNextDateByYears(now, startDate time.Time) string {
	date := startDate

	// Если начальная дата уже в будущем, возвращаем её
	if afterNow(date, now) {
		return date.Format(dateFormat)
	}

	// Начинаем с начальной даты и добавляем годы, пока не превысим now
	yearsToAdd := 1 // Начинаем с +1 года
	for {
		// Создаем новую дату с добавленным годом
		newDate := time.Date(startDate.Year()+yearsToAdd, startDate.Month(), startDate.Day(),
			0, 0, 0, 0, startDate.Location())

		// Обработка 29 февраля
		if newDate.Month() == time.February && newDate.Day() == 29 {
			if !isLeapYear(newDate.Year()) {
				// Переносим на 1 марта
				newDate = time.Date(newDate.Year(), time.March, 1, 0, 0, 0, 0, newDate.Location())
			}
		}

		date = newDate

		// Если дата стала больше now, выходим
		if afterNow(date, now) {
			break
		}

		yearsToAdd++
	}

	return date.Format(dateFormat)
}

// afterNow проверяет, что дата находится после now (без учёта времени)
func afterNow(date, now time.Time) bool {
	// Приводим к началу дня для сравнения только дат
	dateDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return dateDay.After(nowDay)
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

// NextDateHandler обработчик для /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Поддерживаем только GET запросы
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Валидация параметров
	if dateStr == "" {
		http.Error(w, `{"error": "параметр date обязателен"}`, http.StatusBadRequest)
		return
	}

	if repeatStr == "" {
		http.Error(w, `{"error": "параметр repeat обязателен"}`, http.StatusBadRequest)
		return
	}

	// Определяем текущее время
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, `{"error": "неверный формат параметра now"}`, http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату
	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}
