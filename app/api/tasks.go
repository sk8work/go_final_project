package api

import (
	"net/http"

	"github.com/sk8work/go_final_project/app/db"
)

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
	Error string    `json:"error,omitempty"`
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

	var tasks []db.Task
	var err error

	// Устанавливаем лимит (например, 50 задач)
	limit := 50

	if search != "" {
		// Поиск задач
		tasks, err = db.SearchTasks(search, limit)
	} else {
		// Получение всех задач
		tasks, err = db.GetTasks(limit)
	}

	if err != nil {
		writeJSONError(w, "ошибка при получении задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	resp := TasksResponse{Tasks: tasks}
	writeJSON(w, resp)
}
