package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
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

	case "m":
		// Правило: m <дни_месяца> [<месяцы>]
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила m: требуется указать дни")
		}
		return calculateNextDateByMonths(now, startDate, parts[1:])

	case "w":
		// Правило: w <дни_недели>
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила w: требуется указать дни недели")
		}
		return calculateNextDateByWeeks(now, startDate, parts[1])

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// calculateNextDateByDays рассчитывает следующую дату для правила "d"
func calculateNextDateByDays(now, startDate time.Time, interval int) string {
	date := startDate

	// Всегда добавляем хотя бы один интервал
	date = date.AddDate(0, 0, interval)

	// Продолжаем добавлять, пока дата не станет после now
	for !AfterNow(date, now) {
		date = date.AddDate(0, 0, interval)
	}

	return date.Format("20060102")
}

// calculateNextDateByYears рассчитывает следующую дату для правила "y"
func calculateNextDateByYears(now, startDate time.Time) string {
	date := startDate

	// Всегда добавляем хотя бы один год
	date = date.AddDate(1, 0, 0)

	// Обработка 29 февраля
	if date.Month() == time.February && date.Day() == 29 {
		if !isLeapYear(date.Year()) {
			// Переносим на 1 марта
			date = time.Date(date.Year(), time.March, 1, 0, 0, 0, 0, date.Location())
		}
	}

	// Продолжаем добавлять, пока дата не станет после now
	for !AfterNow(date, now) {
		date = date.AddDate(1, 0, 0)

		// Обработка 29 февраля
		if date.Month() == time.February && date.Day() == 29 {
			if !isLeapYear(date.Year()) {
				// Переносим на 1 марта
				date = time.Date(date.Year(), time.March, 1, 0, 0, 0, 0, date.Location())
			}
		}
	}

	return date.Format("20060102")
}

// calculateNextDateByMonths рассчитывает следующую дату для правила "m"
func calculateNextDateByMonths(now, startDate time.Time, params []string) (string, error) {
	// Парсим дни месяца
	daysStr := strings.TrimSpace(params[0])
	if daysStr == "" {
		return "", errors.New("не указаны дни месяца")
	}

	// Разбиваем дни по запятым
	dayItems := strings.Split(daysStr, ",")
	var days []int
	for _, dayStr := range dayItems {
		dayStr = strings.TrimSpace(dayStr)
		if dayStr == "" {
			continue
		}
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", fmt.Errorf("неверный формат дня: %s", dayStr)
		}
		days = append(days, day)
	}

	if len(days) == 0 {
		return "", errors.New("не указаны дни месяца")
	}

	// Парсим месяцы, если указаны
	var months []int
	if len(params) > 1 {
		monthsStr := strings.TrimSpace(params[1])
		monthItems := strings.Split(monthsStr, ",")
		for _, monthStr := range monthItems {
			monthStr = strings.TrimSpace(monthStr)
			if monthStr == "" {
				continue
			}
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				return "", fmt.Errorf("неверный формат месяца: %s", monthStr)
			}
			if month < 1 || month > 12 {
				return "", errors.New("месяц должен быть от 1 до 12")
			}
			months = append(months, month)
		}
	}

	// Если месяцы не указаны, используем все месяцы
	if len(months) == 0 {
		for i := 1; i <= 12; i++ {
			months = append(months, i)
		}
	}

	// Сортируем дни и месяцы для правильного поиска
	sort.Ints(days)
	sort.Ints(months)

	// Находим следующую дату
	return findNextMonthlyDate(now, startDate, days, months)
}

// findNextMonthlyDate находит следующую дату для месячного правила
func findNextMonthlyDate(now, startDate time.Time, days, months []int) (string, error) {
	// Начинаем поиск с дня после now
	current := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Ищем в пределах 5 лет
	for yearOffset := 0; yearOffset < 5; yearOffset++ {
		year := current.Year() + yearOffset

		// Определяем с какого месяца начинать поиск
		startMonthIdx := 0
		if yearOffset == 0 {
			// В текущем году начинаем с текущего месяца
			for i, m := range months {
				if m >= int(current.Month()) {
					startMonthIdx = i
					break
				}
			}
		}

		// Проверяем месяцы
		for monthIdx := startMonthIdx; monthIdx < len(months); monthIdx++ {
			month := months[monthIdx]

			// Определяем с какого дня начинать поиск
			startDayIdx := 0
			if yearOffset == 0 && monthIdx == startMonthIdx && month == int(current.Month()) {
				// В текущем месяце начинаем с текущего дня
				for i, d := range days {
					actualDay := getActualDay(year, month, d)
					if actualDay > current.Day() {
						startDayIdx = i
						break
					}
				}
			}

			// Проверяем дни
			for dayIdx := startDayIdx; dayIdx < len(days); dayIdx++ {
				day := days[dayIdx]
				actualDay := getActualDay(year, month, day)

				// Проверяем, что день существует в месяце
				if actualDay <= 0 {
					continue
				}

				// Создаем кандидата на дату
				candidate := time.Date(year, time.Month(month), actualDay, 0, 0, 0, 0, now.Location())

				// Проверяем, что дата после startDate
				if candidate.Before(startDate) {
					continue
				}

				// Проверяем, что дата после now
				if candidate.After(now) || candidate.Equal(now) {
					return candidate.Format("20060102"), nil
				}
			}
		}

		// Если не нашли в этом году, сбрасываем startMonthIdx для следующего года
	}

	return "", errors.New("не удалось найти следующую дату")
}

func getActualDay(year, month, day int) int {
	if day > 0 {
		// Проверяем, что день существует в месяце
		daysInMonth := daysInMonth(year, time.Month(month))
		if day > daysInMonth {
			return 0
		}
		return day
	}

	// Для отрицательных дней (с конца месяца)
	daysInMonth := daysInMonth(year, time.Month(month))
	actualDay := daysInMonth + day + 1
	if actualDay < 1 || actualDay > daysInMonth {
		return 0
	}
	return actualDay
}

// calculateNextDateByWeeks рассчитывает следующую дату для правила "w"
func calculateNextDateByWeeks(now, startDate time.Time, daysStr string) (string, error) {
	// Парсим дни недели
	daysItems := strings.Split(daysStr, ",")
	var weekdays []int // 1-7, где 1=понедельник, 7=воскресенье

	for _, dayStr := range daysItems {
		dayStr = strings.TrimSpace(dayStr)
		if dayStr == "" {
			continue
		}
		day, err := strconv.Atoi(dayStr)
		if err != nil {
			return "", fmt.Errorf("неверный формат дня недели: %s", dayStr)
		}
		if day < 1 || day > 7 {
			return "", errors.New("день недели должен быть от 1 до 7")
		}
		weekdays = append(weekdays, day)
	}

	if len(weekdays) == 0 {
		return "", errors.New("не указаны дни недели")
	}

	// Начинаем с текущей даты или startDate, если она в будущем
	current := startDate
	if AfterNow(startDate, now) {
		current = startDate
	} else {
		current = now.AddDate(0, 0, 1) // Начинаем со следующего дня
	}

	// Ищем в пределах разумного количества дней
	for i := 0; i < 365*2; i++ {
		candidate := current.AddDate(0, 0, i)

		// Проверяем, что дата после startDate
		if candidate.Before(startDate) || candidate.Equal(startDate) {
			continue
		}

		// Получаем день недели (1=понедельник, 7=воскресенье)
		weekday := int(candidate.Weekday())
		if weekday == 0 {
			weekday = 7 // Воскресенье
		}

		// Проверяем, совпадает ли с одним из указанных дней
		for _, wd := range weekdays {
			if weekday == wd {
				// Проверяем, что дата после now
				if AfterNow(candidate, now) {
					return candidate.Format("20060102"), nil
				}
				break
			}
		}
	}

	return "", errors.New("не удалось найти следующую дату")
}

// daysInMonth возвращает количество дней в месяце
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
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
