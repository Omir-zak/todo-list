// service/task_service.go
package service

import (
	"errors"
	"fmt"
	"time"

	"todo-list/backend/internal/models"
	"todo-list/backend/internal/repository"
)

// NewTaskService создает новый сервис задач (для совместимости с main.go)
func NewTaskService(repo *repository.Repository) *Service {
	return NewService(repo)
}

// TodoService интерфейс для бизнес-логики задач
type TodoService interface {
	CreateTodo(todo *models.Todo) error
	GetTodoByID(id uint) (*models.Todo, error)
	GetAllTodos() ([]models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint) error
	ToggleTodoStatus(id uint) error
	GetCompletedTodos() ([]models.Todo, error)
	GetPendingTodos() ([]models.Todo, error)
}

// CategoryService интерфейс для бизнес-логики категорий
type CategoryService interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
}

// TaskService interface for HTTP handlers (different from TodoService for Wails)
type TaskService interface {
	CreateTask(req *models.CreateTaskRequest) (*models.Todo, error)
	GetTaskByID(id int) (*models.Todo, error)
	GetAllTasks(filter *models.TaskFilter, sort *models.TaskSort) ([]models.Todo, error)
	UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Todo, error)
	DeleteTask(id int) error
	MarkTaskCompleted(id int, completed bool) error
}

// Service объединяет все сервисы
type Service struct {
	Todo     TodoService
	Category CategoryService
}

// todoService реализация TodoService
type todoService struct {
	repo *repository.Repository
}

// categoryService реализация CategoryService
type categoryService struct {
	repo *repository.Repository
}

// taskService implementation for HTTP handlers
type taskService struct {
	repo *repository.Repository
}

// NewService создает новый экземпляр Service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		Todo:     &todoService{repo: repo},
		Category: &categoryService{repo: repo},
	}
}

// NewTaskService создает новый TaskService instance
func NewTaskServiceHandler(repo *repository.Repository) TaskService {
	return &taskService{repo: repo}
}

// Реализация TodoService
func (s *todoService) CreateTodo(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("название задачи обязательно")
	}

	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	return s.repo.Todo.Create(todo)
}

func (s *todoService) GetTodoByID(id uint) (*models.Todo, error) {
	if id == 0 {
		return nil, errors.New("некорректный ID задачи")
	}
	return s.repo.Todo.GetByID(id)
}

func (s *todoService) GetAllTodos() ([]models.Todo, error) {
	return s.repo.Todo.GetAll()
}

func (s *todoService) UpdateTodo(todo *models.Todo) error {
	if todo.ID == 0 {
		return errors.New("некорректный ID задачи")
	}
	if todo.Title == "" {
		return errors.New("название задачи обязательно")
	}

	todo.UpdatedAt = time.Now()
	return s.repo.Todo.Update(todo)
}

func (s *todoService) DeleteTodo(id uint) error {
	if id == 0 {
		return errors.New("некорректный ID задачи")
	}
	return s.repo.Todo.Delete(id)
}

func (s *todoService) ToggleTodoStatus(id uint) error {
	todo, err := s.repo.Todo.GetByID(id)
	if err != nil {
		return fmt.Errorf("задача не найдена: %w", err)
	}

	todo.Completed = !todo.Completed
	todo.UpdatedAt = time.Now()

	return s.repo.Todo.Update(todo)
}

func (s *todoService) GetCompletedTodos() ([]models.Todo, error) {
	return s.repo.Todo.GetByStatus(true)
}

func (s *todoService) GetPendingTodos() ([]models.Todo, error) {
	return s.repo.Todo.GetByStatus(false)
}

// Реализация CategoryService
func (s *categoryService) CreateCategory(category *models.Category) error {
	if category.Name == "" {
		return errors.New("название категории обязательно")
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	return s.repo.Category.Create(category)
}

func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	if id == 0 {
		return nil, errors.New("некорректный ID категории")
	}
	return s.repo.Category.GetByID(id)
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.Category.GetAll()
}

func (s *categoryService) UpdateCategory(category *models.Category) error {
	if category.ID == 0 {
		return errors.New("некорректный ID категории")
	}
	if category.Name == "" {
		return errors.New("название категории обязательно")
	}

	category.UpdatedAt = time.Now()
	return s.repo.Category.Update(category)
}

func (s *categoryService) DeleteCategory(id uint) error {
	if id == 0 {
		return errors.New("некорректный ID категории")
	}
	return s.repo.Category.Delete(id)
}

// Implementation of TaskService methods
func (s *taskService) CreateTask(req *models.CreateTaskRequest) (*models.Todo, error) {
	if req.Title == "" {
		return nil, errors.New("название задачи обязательно")
	}

	todo := &models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Completed:   false,
		CategoryID:  req.CategoryID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Parse due date if provided
	if req.DueDate != "" {
		if dueDate, err := time.Parse("2006-01-02", req.DueDate); err == nil {
			todo.DueDate = &dueDate
		}
	}

	err := s.repo.Todo.Create(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *taskService) GetTaskByID(id int) (*models.Todo, error) {
	if id <= 0 {
		return nil, errors.New("некорректный ID задачи")
	}
	return s.repo.Todo.GetByID(uint(id))
}

func (s *taskService) GetAllTasks(filter *models.TaskFilter, sort *models.TaskSort) ([]models.Todo, error) {
	// For now, return all todos - can be enhanced with filtering and sorting later
	return s.repo.Todo.GetAll()
}

func (s *taskService) UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Todo, error) {
	if id <= 0 {
		return nil, errors.New("некорректный ID задачи")
	}

	todo, err := s.repo.Todo.GetByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("задача не найдена: %w", err)
	}

	// Update fields if provided
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.CategoryID != nil {
		todo.CategoryID = req.CategoryID
	}
	if req.DueDate != nil {
		if *req.DueDate != "" {
			if dueDate, err := time.Parse("2006-01-02", *req.DueDate); err == nil {
				todo.DueDate = &dueDate
			}
		} else {
			todo.DueDate = nil
		}
	}

	todo.UpdatedAt = time.Now()

	err = s.repo.Todo.Update(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *taskService) DeleteTask(id int) error {
	if id <= 0 {
		return errors.New("некорректный ID задачи")
	}
	return s.repo.Todo.Delete(uint(id))
}

func (s *taskService) MarkTaskCompleted(id int, completed bool) error {
	if id <= 0 {
		return errors.New("некорректный ID задачи")
	}

	todo, err := s.repo.Todo.GetByID(uint(id))
	if err != nil {
		return fmt.Errorf("задача не найдена: %w", err)
	}

	todo.Completed = completed
	todo.UpdatedAt = time.Now()

	return s.repo.Todo.Update(todo)
}
