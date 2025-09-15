package wailsbind

import (
	"context"
	"todo-list/backend/internal/models"
	"todo-list/backend/internal/service"
)

// TaskAPI предоставляет API для работы с задачами в Wails
type TaskAPI struct {
	ctx     context.Context
	service *service.Service
}

// NewTaskAPI создает новый экземпляр TaskAPI
func NewTaskAPI(service *service.Service) *TaskAPI {
	return &TaskAPI{
		service: service,
	}
}

// Startup вызывается при старте приложения
func (a *TaskAPI) Startup(ctx context.Context) {
	a.ctx = ctx
}

// GetAllTodos возвращает все задачи
func (a *TaskAPI) GetAllTodos() ([]models.Todo, error) {
	return a.service.Todo.GetAllTodos()
}

// CreateTodo создает новую задачу
func (a *TaskAPI) CreateTodo(title, description string, priority string) (*models.Todo, error) {
	todo := &models.Todo{
		Title:       title,
		Description: description,
		Priority:    priority,
		Completed:   false,
	}

	err := a.service.Todo.CreateTodo(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// UpdateTodo обновляет задачу
func (a *TaskAPI) UpdateTodo(id uint, title, description string, priority string, completed bool) error {
	todo := &models.Todo{
		ID:          id,
		Title:       title,
		Description: description,
		Priority:    priority,
		Completed:   completed,
	}

	return a.service.Todo.UpdateTodo(todo)
}

// DeleteTodo удаляет задачу
func (a *TaskAPI) DeleteTodo(id uint) error {
	return a.service.Todo.DeleteTodo(id)
}

// ToggleTodoStatus переключает статус задачи
func (a *TaskAPI) ToggleTodoStatus(id uint) error {
	return a.service.Todo.ToggleTodoStatus(id)
}

// GetCompletedTodos возвращает завершенные задачи
func (a *TaskAPI) GetCompletedTodos() ([]models.Todo, error) {
	return a.service.Todo.GetCompletedTodos()
}

// GetPendingTodos возвращает незавершенные задачи
func (a *TaskAPI) GetPendingTodos() ([]models.Todo, error) {
	return a.service.Todo.GetPendingTodos()
}

// GetAllCategories возвращает все категории
func (a *TaskAPI) GetAllCategories() ([]models.Category, error) {
	return a.service.Category.GetAllCategories()
}

// CreateCategory создает новую категорию
func (a *TaskAPI) CreateCategory(name, color string) (*models.Category, error) {
	category := &models.Category{
		Name:  name,
		Color: color,
	}

	err := a.service.Category.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}
