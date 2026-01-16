package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sk8work/go_final_project/app/db"
)

// TaskRequest представляет запрос на создание задачи
type TaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTaskHandler обрабатывает POST запрос на добавление задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON запрос
	var req TaskRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, TaskResponse{Error: "ошибка разбора JSON: " + err.Error()})
		return
	}

	// Проверяем обязательные поля
	if req.Title == "" {
		writeJSON(w, TaskResponse{Error: "не указан заголовок задачи"})
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
		task.Date = now.Format("20060102")
	} else {
		// Проверяем формат даты
		date, err := time.Parse("20060102", req.Date)
		if err != nil {
			writeJSON(w, TaskResponse{Error: "неверный формат даты"})
			return
		}
		task.Date = req.Date

		// Если дата в прошлом
		if !AfterNow(date, now) {
			if req.Repeat == "" {
				// Без правила повторения - используем сегодняшнюю дату
				task.Date = now.Format("20060102")
			} else {
				// С правилом повторения - вычисляем следующую дату
				nextDate, err := NextDate(now, req.Date, req.Repeat)
				if err != nil {
					writeJSON(w, TaskResponse{Error: err.Error()})
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
			writeJSON(w, TaskResponse{Error: err.Error()})
			return
		}
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "ошибка при добавлении в базу данных: " + err.Error()})
		return
	}

	// Возвращаем успешный ответ
	writeJSON(w, TaskResponse{ID: fmt.Sprintf("%d", id)})
}
