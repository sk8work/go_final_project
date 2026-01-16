package api

import (
	"fmt"
	"net/http"

	"github.com/sk8work/go_final_project/app/db"
)

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Tasks []TaskJSON `json:"tasks"`
	Error string     `json:"error,omitempty"`
}

// TaskJSON представляет задачу в формате JSON
type TaskJSON struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// convertTaskToJSON конвертирует задачу из БД в JSON формат
func convertTaskToJSON(task *db.Task) TaskJSON {
	return TaskJSON{
		ID:      fmt.Sprintf("%d", task.ID),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}
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

	if search != "" {
		// Поиск задач
		dbTasks, err = db.SearchTasks(search, TasksLimit)
	} else {
		// Получение всех задач
		dbTasks, err = db.GetTasks(TasksLimit)
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
