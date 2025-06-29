package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusToDo       TaskStatus = "to_do"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

func (ts TaskStatus) IsValid() bool {
	return ts == TaskStatusToDo || ts == TaskStatusInProgress || ts == TaskStatusDone || ts == ""
}

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

func (tp TaskPriority) IsValid() bool {
	return tp == TaskPriorityLow || tp == TaskPriorityMedium || tp == TaskPriorityHigh || tp == ""
}

// @Description Task object with all details
type Task struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// @Description Request body for creating a new task
type CreateTaskRequest struct {
	Title       string       `json:"title" validate:"required"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
}

// @Description Request body for updating a task (all fields optional)
type UpdateTaskRequest struct {
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Status      *TaskStatus   `json:"status,omitempty"`
	Priority    *TaskPriority `json:"priority,omitempty"`
}
