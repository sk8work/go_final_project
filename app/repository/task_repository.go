package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sk8work/go_final_project/app/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// TaskFilter определяет параметры фильтрации задач
type TaskFilter struct {
	ID      *int
	Date    *string
	Search  *string
	Limit   *int
	Offset  *int
	OrderBy string // Например: "date ASC", "date DESC", "title"
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
	filter := TaskFilter{ID: &id}
	tasks, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	return &tasks[0], nil
}

// Find возвращает задачи с применением фильтров
func (r *TaskRepository) Find(ctx context.Context, filter TaskFilter) ([]models.Task, error) {
	var (
		queryBuilder strings.Builder
		args         []interface{}
		conditions   []string
	)

	// Базовый запрос
	queryBuilder.WriteString(`
        SELECT id, date, title, comment, repeat
        FROM scheduler
    `)

	// Добавляем условия фильтрации
	if filter.ID != nil {
		conditions = append(conditions, "id = ?")
		args = append(args, *filter.ID)
	}

	if filter.Date != nil {
		conditions = append(conditions, "date = ?")
		args = append(args, *filter.Date)
	}

	if filter.Search != nil {
		searchPattern := "%" + *filter.Search + "%"
		conditions = append(conditions, "(title LIKE ? OR comment LIKE ?)")
		args = append(args, searchPattern, searchPattern)
	}

	// Объединяем условия
	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	// Сортировка
	orderBy := "date ASC" // значение по умолчанию
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	queryBuilder.WriteString(" ORDER BY ")
	queryBuilder.WriteString(orderBy)

	// Лимит и оффсет
	if filter.Limit != nil {
		queryBuilder.WriteString(" LIMIT ?")
		args = append(args, *filter.Limit)
	}

	if filter.Offset != nil {
		queryBuilder.WriteString(" OFFSET ?")
		args = append(args, *filter.Offset)
	}

	// Выполняем запрос
	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
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

// Следующие методы можно оставить для обратной совместимости,
// или переписать их использование на Find

// GetAll возвращает все задачи, отсортированные по дате
func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	return r.Find(ctx, TaskFilter{
		OrderBy: "date ASC",
	})
}

// GetByDate возвращает задачи на конкретную дату
func (r *TaskRepository) GetByDate(ctx context.Context, date string) ([]models.Task, error) {
	return r.Find(ctx, TaskFilter{
		Date:    &date,
		OrderBy: "id ASC",
	})
}

// Search ищет задачи по заголовку или комментарию
func (r *TaskRepository) Search(ctx context.Context, queryStr string) ([]models.Task, error) {
	return r.Find(ctx, TaskFilter{
		Search:  &queryStr,
		OrderBy: "date ASC, id ASC",
	})
}

// IsNotFoundError проверяет, является ли ошибка ошибкой "не найдено"
func IsNotFoundError(err error) bool {
	return err != nil && errors.Is(err, sql.ErrNoRows)
}
