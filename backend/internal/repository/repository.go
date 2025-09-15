// repository/task_repository.go
package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"todo-list/backend/internal/models"
)

type TaskRepository interface {
	Create(task *models.Task) (*models.Task, error)
	GetByID(id int) (*models.Task, error)
	GetAll(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error)
	Update(id int, updates *models.UpdateTaskRequest) (*models.Task, error)
	Delete(id int) error
	MarkCompleted(id int, completed bool) error
}

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) (*models.Task, error) {
	query := `
		INSERT INTO tasks (title, description, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, is_completed, priority, due_date, created_at, updated_at
	`

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	var createdTask models.Task
	err := r.db.QueryRowx(query,
		task.Title,
		task.Description,
		int(task.Priority),
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).StructScan(&createdTask)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &createdTask, nil
}

func (r *taskRepository) GetByID(id int) (*models.Task, error) {
	query := `SELECT id, title, description, is_completed, priority, due_date, created_at, updated_at FROM tasks WHERE id = $1`

	var task models.Task
	err := r.db.Get(&task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) GetAll(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error) {
	query := "SELECT id, title, description, is_completed, priority, due_date, created_at, updated_at FROM tasks"
	args := []interface{}{}
	whereConditions := []string{}
	argIndex := 1

	// Apply filters
	if filter != nil {
		if filter.IsCompleted != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("is_completed = $%d", argIndex))
			args = append(args, *filter.IsCompleted)
			argIndex++
		}

		if filter.Priority != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("priority = $%d", argIndex))
			args = append(args, int(*filter.Priority))
			argIndex++
		}

		if filter.DateFrom != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("created_at >= $%d", argIndex))
			args = append(args, *filter.DateFrom)
			argIndex++
		}

		if filter.DateTo != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("created_at <= $%d", argIndex))
			args = append(args, *filter.DateTo)
			argIndex++
		}
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Apply sorting
	if sort != nil && sort.Field != "" {
		allowedFields := map[string]bool{
			"created_at": true,
			"due_date":   true,
			"priority":   true,
			"title":      true,
		}

		if allowedFields[sort.Field] {
			order := "ASC"
			if sort.Order == "desc" {
				order = "DESC"
			}
			query += fmt.Sprintf(" ORDER BY %s %s", sort.Field, order)
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	var tasks []*models.Task
	err := r.db.Select(&tasks, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) Update(id int, updates *models.UpdateTaskRequest) (*models.Task, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *updates.Title)
		argIndex++
	}

	if updates.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *updates.Description)
		argIndex++
	}

	if updates.IsCompleted != nil {
		setParts = append(setParts, fmt.Sprintf("is_completed = $%d", argIndex))
		args = append(args, *updates.IsCompleted)
		argIndex++
	}

	if updates.Priority != nil {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, int(*updates.Priority))
		argIndex++
	}

	if updates.DueDate != nil {
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", argIndex))
		args = append(args, *updates.DueDate)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE tasks 
		SET %s 
		WHERE id = $%d 
		RETURNING id, title, description, is_completed, priority, due_date, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	var task models.Task
	err := r.db.QueryRowx(query, args...).StructScan(&task)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return &task, nil
}

func (r *taskRepository) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) MarkCompleted(id int, completed bool) error {
	query := "UPDATE tasks SET is_completed = $1, updated_at = $2 WHERE id = $3"
	result, err := r.db.Exec(query, completed, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark task as completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
