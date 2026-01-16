package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sk8work/go_final_project/app/db"
	"github.com/sk8work/go_final_project/app/models"
)

// AddTaskHandler обрабатывает POST запрос на добавление задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON запрос
	var req TaskRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, "неверный формат JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if req.Title == "" {
		writeJSONError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Подготавливаем задачу для БД
	task := &db.Task{
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// Обрабатываем дату
	now := time.Now()

	// Если дата не указана, используем сегодняшнюю
	if req.Date == "" {
		task.Date = models.Today()
	} else {
		// Проверяем формат даты
		date, err := time.Parse(models.DateFormat, req.Date)
		if err != nil {
			writeJSONError(w, "неверный формат даты", http.StatusBadRequest)
			return
		}
		task.Date = req.Date

		// Если дата в прошлом
		if !AfterNow(date, now) {
			if req.Repeat == "" {
				// Без правила повторения - используем сегодняшнюю дату
				task.Date = models.Today()
			} else {
				// С правилом повторения - вычисляем следующую дату
				nextDate, err := NextDate(now, req.Date, req.Repeat)
				if err != nil {
					writeJSONError(w, err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			}
		}
	}

	// Если есть правило повторения, проверяем его
	if req.Repeat != "" {
		// Используем текущую или исправленную дату для проверки
		checkDate := task.Date
		if req.Date != "" {
			checkDate = req.Date
		}

		_, err := NextDate(now, checkDate, req.Repeat)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(task)
	if err != nil {
		writeJSONError(w, "ошибка при добавлении в базу данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSON(w, http.StatusOK, TaskResponse{ID: fmt.Sprintf("%d", id)})
}
