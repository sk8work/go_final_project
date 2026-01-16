package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sk8work/go_final_project/app/models"
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
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// Find возвращает задачи с фильтрацией
func (r *TaskRepository) Find(ctx context.Context, filter TaskFilter) ([]models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1`
	args := []interface{}{}

	if filter.Date != "" {
		query += " AND date = ?"
		args = append(args, filter.Date)
	}

	if filter.Search != "" {
		query += " AND (title LIKE ? OR comment LIKE ?)"
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	query += " ORDER BY date"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
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

// TaskFilter фильтр для поиска задач
type TaskFilter struct {
	Date   string
	Search string
	Limit  int
}
