package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sk8work/go_final_project/app/db"
	"github.com/sk8work/go_final_project/app/models"
)

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
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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
	writeJSON(w, http.StatusOK, taskJSON)
}

// UpdateTaskHandler обрабатывает PUT запрос на обновление задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPut {
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

	// Обновляем задачу в БД
	err = db.UpdateTask(task)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == fmt.Sprintf("задача с id %d не найдена", id) {
			writeJSONError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeJSONError(w, "ошибка при обновлении задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

// DeleteTaskHandler обрабатывает DELETE запрос на удаление задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodDelete {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

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

	// Удаляем задачу из БД
	err = db.DeleteTask(id)
	if err != nil {
		// Проверяем тип ошибки
		if err.Error() == fmt.Sprintf("задача с id %d не найдена", id) {
			writeJSONError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeJSONError(w, "ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
