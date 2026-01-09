package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверка на пустое правило
	if repeat == "" {
		return "", errors.New("повторение не указано")
	}

	// Парсим начальную дату
	startDate, err := time.Parse("20060102", dstart)
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
	if AfterNow(date, now) {
		return date.Format("20060102")
	}

	// Вычисляем следующую дату, добавляя интервалы
	for !AfterNow(date, now) {
		date = date.AddDate(0, 0, interval)
	}

	return date.Format("20060102")
}

// calculateNextDateByYears вычисляет следующую дату для правила y
func calculateNextDateByYears(now, startDate time.Time) string {
	date := startDate

	// Если начальная дата уже в будущем, возвращаем её
	if AfterNow(date, now) {
		return date.Format("20060102")
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
		if AfterNow(date, now) {
			break
		}

		yearsToAdd++
	}

	return date.Format("20060102")
}

// AfterNow проверяет, что дата находится после now (без учёта времени)
func AfterNow(date, now time.Time) bool {
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
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Валидация параметров
	if dateStr == "" {
		writeJSONError(w, "параметр date обязателен", http.StatusBadRequest)
		return
	}

	if repeatStr == "" {
		writeJSONError(w, "параметр repeat обязателен", http.StatusBadRequest)
		return
	}

	// Определяем текущее время
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			writeJSONError(w, "неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату
	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}
