package wailsbind

import (
	"context"

	"todo-list/backend/internal/models"
	"todo-list/backend/internal/service"
)

type TaskAPI struct {
	svc service.TaskService
}

func NewTaskAPI(svc service.TaskService) *TaskAPI {
	return &TaskAPI{svc: svc}
}

// Ниже экспортируемые в JS методы (фронт можешь не писать)
// Они доступны в окне Wails как window.runtime.Events/Bindings (в зав-ти от шаблона)

func (a *TaskAPI) CreateTask(ctx context.Context, req models.CreateTaskRequest) (interface{}, error) {
	return a.svc.CreateTask(&req)
}

func (a *TaskAPI) GetTask(ctx context.Context, id int64) (interface{}, error) {
	return a.svc.GetTaskByID(int(id))
}

func (a *TaskAPI) GetTasks(ctx context.Context, filter *models.TaskFilter, sort *models.TaskSort) (interface{}, error) {
	return a.svc.GetAllTasks(filter, sort)
}

func (a *TaskAPI) UpdateTask(ctx context.Context, id int64, req models.UpdateTaskRequest) (interface{}, error) {
	return a.svc.UpdateTask(int(id), &req)
}

func (a *TaskAPI) DeleteTask(ctx context.Context, id int64) (interface{}, error) {
	return map[string]string{"status": "ok"}, a.svc.DeleteTask(int(id))
}

func (a *TaskAPI) MarkTaskCompleted(ctx context.Context, id int64, completed bool) (interface{}, error) {
	return map[string]string{"status": "ok"}, a.svc.MarkTaskCompleted(int(id), completed)
}
