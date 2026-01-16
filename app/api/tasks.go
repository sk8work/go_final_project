package api

import (
	"net/http"

	"github.com/sk8work/go_final_project/app/db"
)

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Tasks []TaskJSON `json:"tasks"`
	Error string     `json:"error,omitempty"`
}

// TasksHandler обрабатывает GET запрос на получение списка задач
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		writeJSONError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр поиска
	search := r.URL.Query().Get("search")

	var dbTasks []db.Task
	var err error

	// Устанавливаем лимит (например, 50 задач)
	limit := 50

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
	writeJSON(w, resp)
}
