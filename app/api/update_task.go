package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sk8work/go_final_project/app/db"
)

// UpdateTaskHandler обрабатывает PUT запрос на обновление задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPut {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON запрос
	var req struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, "ошибка разбора JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if req.ID == "" {
		writeJSONError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		writeJSONError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в число
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		writeJSONError(w, "неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	// Проверяем существование задачи
	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJSONError(w, "ошибка при получении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if task == nil {
		writeJSONError(w, "задача не найдена", http.StatusNotFound)
		return
	}

	// Подготавливаем обновленную задачу
	now := time.Now()
	updatedTask := &db.Task{
		ID:      id,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// Обрабатываем дату
	if req.Date == "" {
		// Если дата не указана, используем текущую дату задачи или сегодняшнюю
		if task.Date != "" {
			updatedTask.Date = task.Date
		} else {
			updatedTask.Date = now.Format(DateFormat)
		}
	} else {
		// Проверяем формат даты
		date, err := time.Parse(DateFormat, req.Date)
		if err != nil {
			writeJSONError(w, "неверный формат даты", http.StatusBadRequest)
			return
		}
		updatedTask.Date = req.Date

		// Если дата в прошлом
		if !AfterNow(date, now) {
			if req.Repeat == "" {
				// Без правила повторения - используем сегодняшнюю дату
				updatedTask.Date = now.Format(DateFormat)
			} else {
				// С правилом повторения - вычисляем следующую дату
				nextDate, err := NextDate(now, req.Date, req.Repeat)
				if err != nil {
					writeJSONError(w, err.Error(), http.StatusBadRequest)
					return
				}
				updatedTask.Date = nextDate
			}
		}
	}

	// Если есть правило повторения, проверяем его
	if req.Repeat != "" {
		checkDate := updatedTask.Date
		if req.Date != "" {
			checkDate = req.Date
		}

		_, err := NextDate(now, checkDate, req.Repeat)
		if err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Обновляем задачу в БД
	err = db.UpdateTask(updatedTask)
	if err != nil {
		writeJSONError(w, "ошибка при обновлении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
