package api

import (
	"net/http"

	"github.com/sk8work/go_final_project/app/db"
)

const (
	// DefaultTasksLimit лимит задач по умолчанию
	DefaultTasksLimit = 50
	// MaxTasksLimit максимальный лимит задач
	MaxTasksLimit = 1000
)

// TasksHandler обрабатывает GET запрос на получение списка задач
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр поиска
	search := r.URL.Query().Get("search")

	// Получаем лимит из параметров запроса
	limit := GetLimitFromRequest(r)

	var dbTasks []db.Task
	var err error

	if search != "" {
		// Поиск задач
		dbTasks, err = db.SearchTasks(search, limit)
	} else {
		// Получение всех задач
		dbTasks, err = db.GetTasks(limit)
	}

	if err != nil {
		writeJSONError(w, "ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Конвертируем задачи в JSON формат
	tasks := make([]TaskJSON, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, convertTaskToJSON(&dbTask))
	}

	// Возвращаем успешный ответ
	resp := TasksResponse{Tasks: tasks}
	writeJSON(w, http.StatusOK, resp)
}
