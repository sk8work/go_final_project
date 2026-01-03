package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sk8work/go_final_project/app/models"
	_ "strings"
	_ "time"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create создаёт новую задачу
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) (int64, error) {
	query := `
        INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (?, ?, ?, ?)
    `

	result, err := r.db.ExecContext(ctx, query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

// GetByID возвращает задачу по ID
func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE id = ?
    `

	row := r.db.QueryRowContext(ctx, query, id)

	var task models.Task
	err := row.Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// GetAll возвращает все задачи, отсортированные по дате
func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        ORDER BY date
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}

// Update обновляет задачу
func (r *TaskRepository) Update(ctx context.Context, id int, task *models.Task) error {
	query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?
    `

	result, err := r.db.ExecContext(ctx, query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

// Delete удаляет задачу
func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

// GetByDate возвращает задачи на конкретную дату
func (r *TaskRepository) GetByDate(ctx context.Context, date string) ([]models.Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE date = ?
        ORDER BY id
    `

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks by date: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}

// Search ищет задачи по заголовку или комментарию
func (r *TaskRepository) Search(ctx context.Context, queryStr string) ([]models.Task, error) {
	query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        WHERE title LIKE ? OR comment LIKE ?
        ORDER BY date, id
    `

	searchPattern := "%" + queryStr + "%"
	rows, err := r.db.QueryContext(ctx, query, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}
