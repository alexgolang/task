package service

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alexgolang/ishare-task/internal/app/db/sqlite"
	"github.com/alexgolang/ishare-task/internal/app/domain"
	"github.com/google/uuid"
)

func TestTaskService_CreateTask_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDB := ":memory:"
	
	db, err := sqlite.NewDatabase(tempDB)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	logger := log.New(os.Stderr, "INTEGRATION_TEST: ", log.LstdFlags)
	service := NewTaskService(logger, db)

	t.Run("create task successfully", func(t *testing.T) {
		req := &domain.CreateTaskRequest{
			Title:       "Integration Test Task",
			Description: "This is a test task created during integration testing",
			Status:      domain.TaskStatusToDo,
			Priority:    domain.TaskPriorityHigh,
		}

		taskID, err := service.CreateTask(context.Background(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if taskID == uuid.Nil {
			t.Error("Expected a valid UUID, got nil UUID")
		}

		retrievedTask, err := service.GetTask(context.Background(), taskID.String())
		if err != nil {
			t.Fatalf("Failed to retrieve created task: %v", err)
		}

		if retrievedTask.ID != taskID {
			t.Errorf("Expected ID %v, got %v", taskID, retrievedTask.ID)
		}

		if retrievedTask.Title != req.Title {
			t.Errorf("Expected title %v, got %v", req.Title, retrievedTask.Title)
		}

		if retrievedTask.Description != req.Description {
			t.Errorf("Expected description %v, got %v", req.Description, retrievedTask.Description)
		}

		if retrievedTask.Status != req.Status {
			t.Errorf("Expected status %v, got %v", req.Status, retrievedTask.Status)
		}

		if retrievedTask.Priority != req.Priority {
			t.Errorf("Expected priority %v, got %v", req.Priority, retrievedTask.Priority)
		}

		if retrievedTask.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if retrievedTask.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}

		timeDiff := retrievedTask.UpdatedAt.Sub(retrievedTask.CreatedAt)
		if timeDiff > time.Second {
			t.Errorf("Expected CreatedAt and UpdatedAt to be close, got difference: %v", timeDiff)
		}
	})

	t.Run("create task with defaults", func(t *testing.T) {
		req := &domain.CreateTaskRequest{
			Title: "Minimal Task",
		}

		taskID, err := service.CreateTask(context.Background(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		retrievedTask, err := service.GetTask(context.Background(), taskID.String())
		if err != nil {
			t.Fatalf("Failed to retrieve created task: %v", err)
		}

		if retrievedTask.Status != domain.TaskStatusToDo {
			t.Errorf("Expected default status %v, got %v", domain.TaskStatusToDo, retrievedTask.Status)
		}

		if retrievedTask.Priority != domain.TaskPriorityLow {
			t.Errorf("Expected default priority %v, got %v", domain.TaskPriorityLow, retrievedTask.Priority)
		}

		if retrievedTask.Description != "" {
			t.Errorf("Expected empty description, got %v", retrievedTask.Description)
		}
	})

	t.Run("create task with empty title fails", func(t *testing.T) {
		req := &domain.CreateTaskRequest{
			Title:       "",
			Description: "This should fail",
		}

		_, err := service.CreateTask(context.Background(), req)

		if err == nil {
			t.Fatal("Expected error for empty title, got nil")
		}

		expectedError := "create task: title is required"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("create task with invalid status fails", func(t *testing.T) {
		req := &domain.CreateTaskRequest{
			Title:  "Valid Title",
			Status: domain.TaskStatus("invalid_status"),
		}

		_, err := service.CreateTask(context.Background(), req)

		if err == nil {
			t.Fatal("Expected error for invalid status, got nil")
		}

		expectedError := "create task: invalid status"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("create task with invalid priority fails", func(t *testing.T) {
		req := &domain.CreateTaskRequest{
			Title:    "Valid Title",
			Priority: domain.TaskPriority("invalid_priority"),
		}

		_, err := service.CreateTask(context.Background(), req)

		if err == nil {
			t.Fatal("Expected error for invalid priority, got nil")
		}

		expectedError := "create task: invalid priority"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("create multiple tasks", func(t *testing.T) {
		tasks := []*domain.CreateTaskRequest{
			{Title: "Task 1", Priority: domain.TaskPriorityLow},
			{Title: "Task 2", Priority: domain.TaskPriorityMedium},
			{Title: "Task 3", Priority: domain.TaskPriorityHigh},
		}

		var createdIDs []uuid.UUID

		for _, task := range tasks {
			taskID, err := service.CreateTask(context.Background(), task)
			if err != nil {
				t.Fatalf("Failed to create task %v: %v", task.Title, err)
			}
			createdIDs = append(createdIDs, taskID)
		}

		if len(createdIDs) != len(tasks) {
			t.Errorf("Expected %d task IDs, got %d", len(tasks), len(createdIDs))
		}

		idMap := make(map[uuid.UUID]bool)
		for _, id := range createdIDs {
			if idMap[id] {
				t.Errorf("Duplicate task ID found: %v", id)
			}
			idMap[id] = true
		}

		for i, id := range createdIDs {
			retrievedTask, err := service.GetTask(context.Background(), id.String())
			if err != nil {
				t.Errorf("Failed to retrieve task %v: %v", id, err)
				continue
			}

			if retrievedTask.Title != tasks[i].Title {
				t.Errorf("Task %d: expected title %v, got %v", i, tasks[i].Title, retrievedTask.Title)
			}
		}
	})
}