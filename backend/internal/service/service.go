// service/task_service.go
package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"todo-list/backend/internal/models"
	"todo-list/backend/internal/repository"
)

type TaskService interface {
	CreateTask(req *models.CreateTaskRequest) (*models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error)
	UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(id int) error
	MarkTaskCompleted(id int, completed bool) error
}

type taskService struct {
	repo      repository.TaskRepository
	validator *validator.Validate
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo:      repo,
		validator: validator.New(),
	}
}

func (s *taskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		IsCompleted: false,
	}

	return s.repo.Create(task)
}

func (s *taskService) GetTaskByID(id int) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID")
	}

	return s.repo.GetByID(id)
}

func (s *taskService) GetAllTasks(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error) {
	return s.repo.GetAll(filter, sort)
}

func (s *taskService) UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID")
	}

	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if task exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(id, req)
}

func (s *taskService) DeleteTask(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid task ID")
	}

	return s.repo.Delete(id)
}

func (s *taskService) MarkTaskCompleted(id int, completed bool) error {
	if id <= 0 {
		return fmt.Errorf("invalid task ID")
	}

	return s.repo.MarkCompleted(id, completed)
}
