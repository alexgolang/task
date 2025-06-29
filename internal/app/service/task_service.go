package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alexgolang/ishare-task/internal/app/db/sqlite"
	"github.com/alexgolang/ishare-task/internal/app/db/sqlite/sqlc"
	"github.com/alexgolang/ishare-task/internal/app/domain"
	"github.com/google/uuid"
)

type TaskService struct {
	logger *log.Logger
	db     *sqlite.Database
}

func NewTaskService(logger *log.Logger, db *sqlite.Database) *TaskService {
	return &TaskService{
		logger: logger,
		db:     db,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.CreateTaskRequest) (uuid.UUID, error) {
	if task.Title == "" {
		return uuid.UUID{}, fmt.Errorf("create task: title is required")
	}

	status := domain.TaskStatusToDo
	priority := domain.TaskPriorityLow

	if task.Status != "" {
		status = task.Status
	}

	if task.Priority != "" {
		priority = task.Priority
	}

	if !status.IsValid() {
		return uuid.UUID{}, fmt.Errorf("create task: invalid status")
	}

	if !priority.IsValid() {
		return uuid.UUID{}, fmt.Errorf("create task: invalid priority")
	}

	id := uuid.New()
	err := s.db.Queries.CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          id.String(),
		Title:       task.Title,
		Description: sql.NullString{String: task.Description, Valid: task.Description != ""},
		Status:      status,
		Priority:    priority,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create task: %w", err)
	}

	return id, nil
}

func (s *TaskService) GetTasks(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.db.Queries.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	domainTasks := make([]domain.Task, len(tasks))
	for i, task := range tasks {
		domainTasks[i] = toDomain(task)
	}

	return domainTasks, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (domain.Task, error) {
	if id == "" {
		return domain.Task{}, fmt.Errorf("get task: id is required")
	}

	task, err := s.db.Queries.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	return toDomain(task), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, task *domain.UpdateTaskRequest) error {
	if id == "" {
		return fmt.Errorf("update task: id is required")
	}

	if task.Status != nil && !task.Status.IsValid() {
		return fmt.Errorf("update task: invalid status")
	}

	if task.Priority != nil && !task.Priority.IsValid() {
		return fmt.Errorf("update task: invalid priority")
	}

	title := ""
	if task.Title != nil {
		title = *task.Title
	}

	status := domain.TaskStatus("")
	if task.Status != nil {
		status = *task.Status
	}

	priority := domain.TaskPriority("")
	if task.Priority != nil {
		priority = *task.Priority
	}

	description := ""
	if task.Description != nil {
		description = *task.Description
	}

	err := s.db.Queries.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:          id,
		Title:       title,
		Description: sql.NullString{String: description, Valid: description != ""},
		Status:      status,
		Priority:    priority,
		UpdatedAt:   time.Now().UTC(),
	})

	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("delete task: id is required")
	}

	result, err := s.db.Queries.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if result == 0 {
		return fmt.Errorf("delete task: task not found")
	}

	return nil
}

func toDomain(task sqlc.Task) domain.Task {
	return domain.Task{
		ID:          uuid.MustParse(task.ID),
		Title:       task.Title,
		Description: task.Description.String,
		Status:      task.Status,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
