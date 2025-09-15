// handler/task_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"todo-list/backend/internal/models"
	"todo-list/backend/internal/service"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.service.CreateTask(&req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, task)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	filter := h.parseFilter(r)
	sort := h.parseSort(r)

	tasks, err := h.service.GetAllTasks(filter, sort)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.service.UpdateTask(id, &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func (h *TaskHandler) MarkTaskCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req struct {
		Completed bool `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.service.MarkTaskCompleted(id, req.Completed)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	h.writeSuccess(w, http.StatusOK, map[string]string{"message": "Task status updated successfully"})
}

func (h *TaskHandler) parseFilter(r *http.Request) *models.TaskFilter {
	query := r.URL.Query()
	filter := &models.TaskFilter{}

	if completedStr := query.Get("completed"); completedStr != "" {
		if completed, err := strconv.ParseBool(completedStr); err == nil {
			filter.IsCompleted = &completed
		}
	}

	if priorityStr := query.Get("priority"); priorityStr != "" {
		switch priorityStr {
		case "low":
			priority := models.Low
			filter.Priority = &priority
		case "medium":
			priority := models.Medium
			filter.Priority = &priority
		case "high":
			priority := models.High
			filter.Priority = &priority
		}
	}

	if dateFromStr := query.Get("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filter.DateFrom = &dateFrom
		}
	}

	if dateToStr := query.Get("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			filter.DateTo = &dateTo
		}
	}

	return filter
}

func (h *TaskHandler) parseSort(r *http.Request) *models.TaskSort {
	query := r.URL.Query()
	return &models.TaskSort{
		Field: query.Get("sort_by"),
		Order: query.Get("sort_order"),
	}
}

func (h *TaskHandler) writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

func (h *TaskHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}
