package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexgolang/ishare-task/internal/app/common/server"
	"github.com/alexgolang/ishare-task/internal/app/domain"
	"github.com/alexgolang/ishare-task/internal/app/service"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with title, description, status, and priority
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body domain.CreateTaskRequest true "Task data"
// @Success 201 {object} map[string]string "Task created successfully"
// @Failure 400 {object} server.ErrorResponse "Bad request"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RespondError(err, w, r)
		return
	}

	id, err := h.taskService.CreateTask(r.Context(), &domain.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	})

	if err != nil {
		server.RespondError(err, w, r)
		return
	}

	server.RespondOK(fmt.Sprintf("Task %s created", id), w, r)
}

// GetTask godoc
// @Summary Get task by ID
// @Description Get a single task by its UUID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID (UUID)"
// @Success 200 {object} domain.Task "Task found"
// @Failure 404 {object} server.ErrorResponse "Task not found"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.taskService.GetTask(r.Context(), id)

	if err != nil {
		if err == sql.ErrNoRows {
			server.RespondNotFound("Task not found", w, r)
			return
		}
		server.RespondError(err, w, r)
		return
	}

	server.RespondOK(task, w, r)
}

// UpdateTask godoc
// @Summary Update task
// @Description Update a task's fields (partial update supported)
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID (UUID)"
// @Param task body domain.UpdateTaskRequest true "Task update data"
// @Success 200 {object} map[string]string "Task updated successfully"
// @Failure 400 {object} server.ErrorResponse "Bad request"
// @Failure 404 {object} server.ErrorResponse "Task not found"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.RespondError(err, w, r)
		return
	}

	err := h.taskService.UpdateTask(r.Context(), id, &domain.UpdateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	})

	if err != nil {
		server.RespondError(err, w, r)
		return
	}

	server.RespondOK(fmt.Sprintf("Task %s updated", id), w, r)
}

// DeleteTask godoc
// @Summary Delete task
// @Description Delete a task by its UUID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID (UUID)"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 404 {object} server.ErrorResponse "Task not found"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.taskService.DeleteTask(r.Context(), id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			server.RespondNotFound(err.Error(), w, r)
			return
		}
		server.RespondError(err, w, r)
		return
	}

	server.RespondOK(fmt.Sprintf("Task %s deleted", id), w, r)
}


// ListTasks godoc
// @Summary List all tasks
// @Description Get all tasks in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} domain.Task "List of tasks"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetTasks(r.Context())

	if err != nil {
		server.RespondError(err, w, r)
		return
	}

	server.RespondOK(tasks, w, r)
}
