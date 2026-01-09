package db

import (
	"database/sql"
	"fmt"
)

// Task представляет задачу в планировщике
type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`    // Формат: YYYYMMDD
	Title   string `json:"title"`   // Заголовок задачи
	Comment string `json:"comment"` // Комментарий
	Repeat  string `json:"repeat"`  // Правило повторения
}

// AddTask добавляет задачу в базу данных
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении задачи: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении ID: %w", err)
	}

	return id, nil
}

// GetTaskByID возвращает задачу по ID
func GetTaskByID(id int64) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	row := DB.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Задача не найдена
		}
		return nil, fmt.Errorf("ошибка при чтении задачи: %w", err)
	}

	return &task, nil
}
