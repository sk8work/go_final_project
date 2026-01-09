package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sk8work/go_final_project/app/db"
)

// UpdateTaskRequest представляет запрос на обновление задачи
type UpdateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// convertTaskToJSON конвертирует задачу из БД в JSON формат
func convertTaskToJSON(task *db.Task) TaskJSON {
	return TaskJSON{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
}

// GetTaskHandler обрабатывает GET запрос на получение задачи
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, "неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJSONError(w, "ошибка при получении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if task == nil {
		writeJSONError(w, "задача не найдена", http.StatusNotFound)
		return
	}

	// Конвертируем задачу в JSON формат и возвращаем
	taskJSON := convertTaskToJSON(task)
	writeJSON(w, taskJSON)
}

// UpdateTaskHandler обрабатывает PUT запрос на обновление задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON запрос
	var req UpdateTaskRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, TaskResponse{Error: "ошибка разбора JSON: " + err.Error()})
		return
	}

	// Проверяем обязательные поля
	if req.ID == "" {
		writeJSON(w, TaskResponse{Error: "не указан идентификатор задачи"})
		return
	}

	if req.Title == "" {
		writeJSON(w, TaskResponse{Error: "не указан заголовок задачи"})
		return
	}

	// Преобразуем ID в число
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "неверный формат идентификатора"})
		return
	}

	// Подготавливаем задачу для БД
	task := &db.Task{
		ID:      id,
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

	// Обновляем задачу в БД
	err = db.UpdateTask(task)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "ошибка при обновлении задачи: " + err.Error()})
		return
	}

	// Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, map[string]interface{}{})
}
