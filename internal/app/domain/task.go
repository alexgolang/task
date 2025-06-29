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
	ID          uuid.UUID
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// @Description Request body for creating a new task
type CreateTaskRequest struct {
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
}

// @Description Request body for updating a task (all fields optional)
type UpdateTaskRequest struct {
	Title       *string
	Description *string
	Status      *TaskStatus
	Priority    *TaskPriority
}
