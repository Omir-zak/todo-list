package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Task структура задачи
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority"` // low, medium, high
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}

// TaskManager управляет задачами
type TaskManager struct {
	tasks    []Task
	nextID   int
	filename string
}

// NewTaskManager создает новый менеджер задач
func NewTaskManager() *TaskManager {
	// Получаем путь к файлу данных
	homeDir, _ := os.UserHomeDir()
	filename := filepath.Join(homeDir, ".todo-list.json")

	tm := &TaskManager{
		tasks:    []Task{},
		nextID:   1,
		filename: filename,
	}

	// Загружаем существующие задачи
	tm.loadTasks()

	return tm
}

// GetTasks возвращает все задачи
func (a *App) GetTasks() []Task {
	if a.taskManager == nil {
		a.taskManager = NewTaskManager()
	}
	return a.taskManager.tasks
}

// AddTask добавляет новую задачу
func (a *App) AddTask(title, description, priority string, dueDate string) Task {
	if a.taskManager == nil {
		a.taskManager = NewTaskManager()
	}

	if title == "" {
		return Task{} // Валидация на пустой ввод
	}

	var due time.Time
	if dueDate != "" {
		due, _ = time.Parse("2006-01-02T15:04", dueDate)
	}

	task := Task{
		ID:          a.taskManager.nextID,
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     due,
		CreatedAt:   time.Now(),
		Completed:   false,
	}

	a.taskManager.tasks = append(a.taskManager.tasks, task)
	a.taskManager.nextID++

	// Сохраняем изменения
	a.taskManager.saveTasks()

	return task
}

// DeleteTask удаляет задачу по ID
func (a *App) DeleteTask(id int) bool {
	if a.taskManager == nil {
		return false
	}

	for i, task := range a.taskManager.tasks {
		if task.ID == id {
			a.taskManager.tasks = append(a.taskManager.tasks[:i], a.taskManager.tasks[i+1:]...)
			a.taskManager.saveTasks()
			return true
		}
	}
	return false
}

// ToggleTask переключает состояние выполнения задачи
func (a *App) ToggleTask(id int) bool {
	if a.taskManager == nil {
		return false
	}

	for i, task := range a.taskManager.tasks {
		if task.ID == id {
			a.taskManager.tasks[i].Completed = !task.Completed
			a.taskManager.saveTasks()
			return true
		}
	}
	return false
}

// GetFilteredTasks возвращает отфильтрованные задачи
func (a *App) GetFilteredTasks(filter string) []Task {
	if a.taskManager == nil {
		a.taskManager = NewTaskManager()
	}

	var filtered []Task

	for _, task := range a.taskManager.tasks {
		switch filter {
		case "active":
			if !task.Completed {
				filtered = append(filtered, task)
			}
		case "completed":
			if task.Completed {
				filtered = append(filtered, task)
			}
		default: // "all"
			filtered = append(filtered, task)
		}
	}

	return filtered
}

// loadTasks загружает задачи из файла
func (tm *TaskManager) loadTasks() {
	data, err := os.ReadFile(tm.filename)
	if err != nil {
		return // Файл не существует или ошибка чтения
	}

	var savedData struct {
		Tasks  []Task `json:"tasks"`
		NextID int    `json:"next_id"`
	}

	if err := json.Unmarshal(data, &savedData); err != nil {
		return
	}

	tm.tasks = savedData.Tasks
	tm.nextID = savedData.NextID
}

// saveTasks сохраняет задачи в файл
func (tm *TaskManager) saveTasks() {
	data := struct {
		Tasks  []Task `json:"tasks"`
		NextID int    `json:"next_id"`
	}{
		Tasks:  tm.tasks,
		NextID: tm.nextID,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling tasks: %v\n", err)
		return
	}

	if err := os.WriteFile(tm.filename, jsonData, 0644); err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
	}
}
