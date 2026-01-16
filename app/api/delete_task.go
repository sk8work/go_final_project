package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sk8work/go_final_project/app/db"
)

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

	// Удаляем задачу из БД
	err = db.DeleteTask(id)
	if err != nil {
		writeJSONError(w, "ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
