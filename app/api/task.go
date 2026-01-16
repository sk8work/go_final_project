package api

import (
	"github.com/sk8work/go_final_project/app/db"
	"net/http"
	"strconv"
)

// UpdateTaskRequest представляет запрос на обновление задачи
type UpdateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
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
