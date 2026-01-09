package db

import (
	"database/sql"
	"fmt"
	"time"
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

// UpdateTask обновляет задачу в базе данных
func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с id %d не найдена", task.ID)
	}

	return nil
}

// GetTasks возвращает список задач с лимитом
func GetTasks(limit int) ([]Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе задач: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// SearchTasks ищет задачи по заголовку или комментарию
func SearchTasks(search string, limit int) ([]Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	// Пытаемся парсить дату в формате DD.MM.YYYY
	if date, err := time.Parse("02.01.2006", search); err == nil {
		// Если search - это дата, ищем по дате
		dateStr := date.Format("20060102")
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`

		rows, err := DB.Query(query, dateStr, limit)
		if err != nil {
			return nil, fmt.Errorf("ошибка при поиске задач по дате: %w", err)
		}
		defer rows.Close()

		return scanTasks(rows)
	}

	// Иначе ищем по заголовку или комментарию
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          WHERE title LIKE ? OR comment LIKE ? 
	          ORDER BY date LIMIT ?`

	searchPattern := "%" + search + "%"
	rows, err := DB.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске задач: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks сканирует строки из запроса в массив задач
func scanTasks(rows *sql.Rows) ([]Task, error) {
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании задачи: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по задачам: %w", err)
	}

	// Возвращаем пустой слайс вместо nil, чтобы в JSON было "tasks": []
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}
