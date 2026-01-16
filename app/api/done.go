package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sk8work/go_final_project/app/db"
)

// DoneTaskHandler обрабатывает POST запрос на завершение задачи
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
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

	// Если задача без повторения - удаляем её
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSONError(w, "ошибка при удалении задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Если задача с повторением - вычисляем следующую дату
		now := time.Now()

		// Используем текущую дату задачи как начальную для расчета
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "ошибка при вычислении следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			writeJSONError(w, "ошибка при обновлении задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
