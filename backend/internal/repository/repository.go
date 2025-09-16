// repository/task_repository.go
package repository

import (
	"database/sql"
	"time"
	"todo-list/backend/internal/models"
)

// NewTaskRepository создает новый репозиторий задач (для совместимости с start.go)
func NewTaskRepository(db *sql.DB) *Repository {
	return NewRepository(db)
}

// TodoRepository интерфейс для работы с задачами
type TodoRepository interface {
	Create(todo *models.Todo) error
	GetByID(id uint) (*models.Todo, error)
	GetAll() ([]models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id uint) error
	GetByStatus(completed bool) ([]models.Todo, error)
}

// CategoryRepository интерфейс для работы с категориями
type CategoryRepository interface {
	Create(category *models.Category) error
	GetByID(id uint) (*models.Category, error)
	GetAll() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

// Repository объединяет все репозитории
type Repository struct {
	Todo     TodoRepository
	Category CategoryRepository
}

// todoRepo реализация TodoRepository
type todoRepo struct {
	db *sql.DB
}

// categoryRepo реализация CategoryRepository
type categoryRepo struct {
	db *sql.DB
}

// NewRepository создает новый экземпляр Repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Todo:     &todoRepo{db: db},
		Category: &categoryRepo{db: db},
	}
}

// Реализация TodoRepository

func (r *todoRepo) Create(todo *models.Todo) error {
	query := `
		INSERT INTO todos (title, description, completed, priority, due_date, category_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	return r.db.QueryRow(query, todo.Title, todo.Description, todo.Completed,
		todo.Priority, todo.DueDate, todo.CategoryID,
		todo.CreatedAt, todo.UpdatedAt).Scan(&todo.ID)
}

func (r *todoRepo) GetByID(id uint) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `
		SELECT id, title, description, completed, priority, due_date, 
		       category_id, created_at, updated_at 
		FROM todos WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Completed,
		&todo.Priority, &todo.DueDate, &todo.CategoryID,
		&todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *todoRepo) GetAll() ([]models.Todo, error) {
	query := `
		SELECT id, title, description, completed, priority, due_date, 
		       category_id, created_at, updated_at 
		FROM todos ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Completed,
			&todo.Priority, &todo.DueDate, &todo.CategoryID,
			&todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

func (r *todoRepo) Update(todo *models.Todo) error {
	query := `
		UPDATE todos SET title = $1, description = $2, completed = $3, 
		                 priority = $4, due_date = $5, category_id = $6, 
		                 updated_at = $7 
		WHERE id = $8`

	todo.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, todo.Title, todo.Description, todo.Completed,
		todo.Priority, todo.DueDate, todo.CategoryID,
		todo.UpdatedAt, todo.ID)
	return err
}

func (r *todoRepo) Delete(id uint) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *todoRepo) GetByStatus(completed bool) ([]models.Todo, error) {
	query := `
		SELECT id, title, description, completed, priority, due_date, 
		       category_id, created_at, updated_at 
		FROM todos WHERE completed = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, completed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.Title, &todo.Description, &todo.Completed,
			&todo.Priority, &todo.DueDate, &todo.CategoryID,
			&todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, rows.Err()
}

// Реализация CategoryRepository

func (r *categoryRepo) Create(category *models.Category) error {
	query := `
		INSERT INTO categories (name, color, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	return r.db.QueryRow(query, category.Name, category.Color,
		category.CreatedAt, category.UpdatedAt).Scan(&category.ID)
}

func (r *categoryRepo) GetByID(id uint) (*models.Category, error) {
	category := &models.Category{}
	query := `
		SELECT id, name, color, created_at, updated_at 
		FROM categories WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&category.ID, &category.Name, &category.Color,
		&category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepo) GetAll() ([]models.Category, error) {
	query := `
		SELECT id, name, color, created_at, updated_at 
		FROM categories ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID, &category.Name, &category.Color,
			&category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepo) Update(category *models.Category) error {
	query := `
		UPDATE categories SET name = $1, color = $2, updated_at = $3 
		WHERE id = $4`

	category.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, category.Name, category.Color,
		category.UpdatedAt, category.ID)
	return err
}

func (r *categoryRepo) Delete(id uint) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
