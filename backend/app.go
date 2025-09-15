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

// App структура приложения
type App struct {
	taskManager *TaskManager
}

// NewApp создает новый экземпляр приложения
func NewApp() *App {
	return &App{}
}

// TaskManager управляет задачами
type TaskManager struct {
	tasks    []Task
	nextID   int
	filename string
}

// NewTaskManager создает новый менеджер задач
func NewTaskManager() *TaskManager {
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

// GetTasksByDateFilter возвращает задачи по фильтру даты
func (a *App) GetTasksByDateFilter(filter string) []Task {
	if a.taskManager == nil {
		a.taskManager = NewTaskManager()
	}

	var filtered []Task
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekFromNow := today.AddDate(0, 0, 7)

	for _, task := range a.taskManager.tasks {
		if task.DueDate.IsZero() {
			continue // Пропускаем задачи без даты
		}

		taskDate := time.Date(task.DueDate.Year(), task.DueDate.Month(), task.DueDate.Day(), 0, 0, 0, 0, task.DueDate.Location())

		switch filter {
		case "today":
			if taskDate.Equal(today) {
				filtered = append(filtered, task)
			}
		case "week":
			if taskDate.After(today.AddDate(0, 0, -1)) && taskDate.Before(weekFromNow) {
				filtered = append(filtered, task)
			}
		case "overdue":
			if taskDate.Before(today) && !task.Completed {
				filtered = append(filtered, task)
			}
		default: // "all"
			filtered = append(filtered, task)
		}
	}

	return filtered
}

// GetSortedTasks возвращает отсортированные задачи
func (a *App) GetSortedTasks(sortBy string, ascending bool) []Task {
	if a.taskManager == nil {
		a.taskManager = NewTaskManager()
	}

	tasks := make([]Task, len(a.taskManager.tasks))
	copy(tasks, a.taskManager.tasks)

	switch sortBy {
	case "date":
		if ascending {
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if tasks[i].CreatedAt.After(tasks[j].CreatedAt) {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if tasks[i].CreatedAt.Before(tasks[j].CreatedAt) {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		}
	case "priority":
		priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1, "": 0}
		if ascending {
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if priorityOrder[tasks[i].Priority] > priorityOrder[tasks[j].Priority] {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if priorityOrder[tasks[i].Priority] < priorityOrder[tasks[j].Priority] {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		}
	case "dueDate":
		if ascending {
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if !tasks[i].DueDate.IsZero() && !tasks[j].DueDate.IsZero() && tasks[i].DueDate.After(tasks[j].DueDate) {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		} else {
			for i := 0; i < len(tasks); i++ {
				for j := i + 1; j < len(tasks); j++ {
					if !tasks[i].DueDate.IsZero() && !tasks[j].DueDate.IsZero() && tasks[i].DueDate.Before(tasks[j].DueDate) {
						tasks[i], tasks[j] = tasks[j], tasks[i]
					}
				}
			}
		}
	}

	return tasks
}

// GetCombinedFilteredTasks возвращает задачи с комбинированными фильтрами
func (a *App) GetCombinedFilteredTasks(statusFilter, dateFilter, sortBy string, ascending bool) []Task {
	if a.taskManager == nil {
		a.taskManager = NewTaskManager()
	}

	// Сначала применяем фильтр по статусу
	var filtered []Task
	for _, task := range a.taskManager.tasks {
		switch statusFilter {
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

	// Применяем фильтр по дате
	if dateFilter != "" && dateFilter != "all" {
		var dateFiltered []Task
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		weekFromNow := today.AddDate(0, 0, 7)

		for _, task := range filtered {
			if task.DueDate.IsZero() && dateFilter != "all" {
				continue
			}

			taskDate := time.Date(task.DueDate.Year(), task.DueDate.Month(), task.DueDate.Day(), 0, 0, 0, 0, task.DueDate.Location())

			switch dateFilter {
			case "today":
				if taskDate.Equal(today) {
					dateFiltered = append(dateFiltered, task)
				}
			case "week":
				if taskDate.After(today.AddDate(0, 0, -1)) && taskDate.Before(weekFromNow) {
					dateFiltered = append(dateFiltered, task)
				}
			case "overdue":
				if taskDate.Before(today) && !task.Completed {
					dateFiltered = append(dateFiltered, task)
				}
			}
		}
		filtered = dateFiltered
	}

	// Применяем сортировку
	if sortBy != "" {
		switch sortBy {
		case "date":
			if ascending {
				for i := 0; i < len(filtered); i++ {
					for j := i + 1; j < len(filtered); j++ {
						if filtered[i].CreatedAt.After(filtered[j].CreatedAt) {
							filtered[i], filtered[j] = filtered[j], filtered[i]
						}
					}
				}
			} else {
				for i := 0; i < len(filtered); i++ {
					for j := i + 1; j < len(filtered); j++ {
						if filtered[i].CreatedAt.Before(filtered[j].CreatedAt) {
							filtered[i], filtered[j] = filtered[j], filtered[i]
						}
					}
				}
			}
		case "priority":
			priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1, "": 0}
			if ascending {
				for i := 0; i < len(filtered); i++ {
					for j := i + 1; j < len(filtered); j++ {
						if priorityOrder[filtered[i].Priority] > priorityOrder[filtered[j].Priority] {
							filtered[i], filtered[j] = filtered[j], filtered[i]
						}
					}
				}
			} else {
				for i := 0; i < len(filtered); i++ {
					for j := i + 1; j < len(filtered); j++ {
						if priorityOrder[filtered[i].Priority] < priorityOrder[filtered[j].Priority] {
							filtered[i], filtered[j] = filtered[j], filtered[i]
						}
					}
				}
			}
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
