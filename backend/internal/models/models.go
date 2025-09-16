package models

import (
	"time"
)

// Priority represents task priority levels
type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

// Todo представляет задачу в списке дел
type Todo struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CategoryID  *uint      `json:"category_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Category представляет категорию задач
type Category struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request structs for API handlers
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	CategoryID  *uint  `json:"category_id"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	DueDate     *string `json:"due_date"`
	Completed   *bool   `json:"completed"`
	CategoryID  *uint   `json:"category_id"`
}

// Filter and sort structs
type TaskFilter struct {
	IsCompleted *bool      `json:"is_completed"`
	Priority    *Priority  `json:"priority"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	CategoryID  *uint      `json:"category_id"`
}

type TaskSort struct {
	Field string `json:"field"` // id, title, priority, due_date, created_at
	Order string `json:"order"` // asc, desc
}
