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
		writeJSON(w, TaskResponse{Error: "не указан идентификатор задачи"})
		return
	}

	// Преобразуем ID в число
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "неверный формат идентификатора"})
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "ошибка при получении задачи: " + err.Error()})
		return
	}

	if task == nil {
		writeJSON(w, TaskResponse{Error: "задача не найдена"})
		return
	}

	// Если задача без повторения - удаляем её
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, TaskResponse{Error: "ошибка при удалении задачи: " + err.Error()})
			return
		}
	} else {
		// Если задача с повторением - вычисляем следующую дату
		now := time.Now()

		// Используем текущую дату задачи как начальную для расчета
		// Это важно: каждый раз мы рассчитываем следующую дату от ТЕКУЩЕЙ ДАТЫ ЗАДАЧИ
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, TaskResponse{Error: "ошибка при вычислении следующей даты: " + err.Error()})
			return
		}

		// Обновляем дату задачи
		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			writeJSON(w, TaskResponse{Error: "ошибка при обновлении задачи: " + err.Error()})
			return
		}
	}

	// Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, map[string]interface{}{})
}
