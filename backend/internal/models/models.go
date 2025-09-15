package models

import (
	"time"

	"gorm.io/gorm"
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
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	Priority    string         `json:"priority" gorm:"default:'medium'"` // low, medium, high
	DueDate     *time.Time     `json:"due_date"`
	UserID      *uint          `json:"user_id" gorm:"index"`
	CategoryID  *uint          `json:"category_id" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// User представляет пользователя системы
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Todos     []Todo         `json:"todos" gorm:"foreignKey:UserID"`
}

// Category представляет категорию задач
type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Color     string         `json:"color" gorm:"default:'#007bff'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Todos     []Todo         `json:"todos" gorm:"foreignKey:CategoryID"`
}

// TodoWithRelations структура для получения todo с связанными данными
type TodoWithRelations struct {
	Todo
	User     *User     `json:"user,omitempty"`
	Category *Category `json:"category,omitempty"`
}

// Request structs for API handlers
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	UserID      *uint  `json:"user_id"`
	CategoryID  *uint  `json:"category_id"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	DueDate     *string `json:"due_date"`
	Completed   *bool   `json:"completed"`
	UserID      *uint   `json:"user_id"`
	CategoryID  *uint   `json:"category_id"`
}

// Filter and sort structs
type TaskFilter struct {
	IsCompleted *bool      `json:"is_completed"`
	Priority    *Priority  `json:"priority"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	UserID      *uint      `json:"user_id"`
	CategoryID  *uint      `json:"category_id"`
}

type TaskSort struct {
	Field string `json:"field"` // id, title, priority, due_date, created_at
	Order string `json:"order"` // asc, desc
}
