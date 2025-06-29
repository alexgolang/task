package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alexgolang/ishare-task/internal/app/db/sqlite/sqlc"
	"github.com/alexgolang/ishare-task/internal/app/domain"
	"github.com/alexgolang/ishare-task/internal/app/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

type TestTaskService struct {
	logger  *log.Logger
	querier sqlc.Querier
}

func newTestTaskService(logger *log.Logger, querier sqlc.Querier) *TestTaskService {
	return &TestTaskService{
		logger:  logger,
		querier: querier,
	}
}

func (s *TestTaskService) GetTask(ctx context.Context, id string) (domain.Task, error) {
	if id == "" {
		s.logger.Printf("get task: id is required")
		return domain.Task{}, fmt.Errorf("get task: id is required")
	}

	task, err := s.querier.GetTask(ctx, id)
	if err != nil {
		s.logger.Printf("get task: %v", err)
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	return toDomain(task), nil
}

func TestTaskService_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := log.New(os.Stderr, "TEST: ", log.LstdFlags)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedID := uuid.New()
		expectedTask := sqlc.Task{
			ID:          expectedID.String(),
			Title:       "Test Task",
			Description: sql.NullString{String: "Test Description", Valid: true},
			Status:      domain.TaskStatusToDo,
			Priority:    domain.TaskPriorityHigh,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockQuerier := mocks.NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetTask(gomock.Any(), expectedID.String()).
			Return(expectedTask, nil).
			Times(1)

		service := newTestTaskService(logger, mockQuerier)

		result, err := service.GetTask(context.Background(), expectedID.String())

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.ID != expectedID {
			t.Errorf("Expected ID %v, got %v", expectedID, result.ID)
		}

		if result.Title != expectedTask.Title {
			t.Errorf("Expected title %v, got %v", expectedTask.Title, result.Title)
		}

		if result.Status != expectedTask.Status {
			t.Errorf("Expected status %v, got %v", expectedTask.Status, result.Status)
		}

		if result.Description != expectedTask.Description.String {
			t.Errorf("Expected description %v, got %v", expectedTask.Description.String, result.Description)
		}

		if result.Priority != expectedTask.Priority {
			t.Errorf("Expected priority %v, got %v", expectedTask.Priority, result.Priority)
		}
	})

	t.Run("empty ID validation", func(t *testing.T) {
		mockQuerier := mocks.NewMockQuerier(ctrl)
		service := newTestTaskService(logger, mockQuerier)

		_, err := service.GetTask(context.Background(), "")

		if err == nil {
			t.Fatal("Expected error for empty ID, got nil")
		}

		expectedError := "get task: id is required"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("task not found", func(t *testing.T) {
		taskID := uuid.New().String()
		
		mockQuerier := mocks.NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetTask(gomock.Any(), taskID).
			Return(sqlc.Task{}, sql.ErrNoRows).
			Times(1)

		service := newTestTaskService(logger, mockQuerier)

		_, err := service.GetTask(context.Background(), taskID)

		if err == nil {
			t.Fatal("Expected error for non-existent task, got nil")
		}

		expectedError := "get task: sql: no rows in result set"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("database error", func(t *testing.T) {
		taskID := uuid.New().String()
		dbError := sql.ErrConnDone
		
		mockQuerier := mocks.NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetTask(gomock.Any(), taskID).
			Return(sqlc.Task{}, dbError).
			Times(1)

		service := newTestTaskService(logger, mockQuerier)

		_, err := service.GetTask(context.Background(), taskID)

		if err == nil {
			t.Fatal("Expected database error, got nil")
		}

		expectedError := "get task: sql: connection is already closed"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})

	t.Run("invalid UUID conversion", func(t *testing.T) {
		taskID := "valid-id"
		invalidTask := sqlc.Task{
			ID:          "invalid-uuid",
			Title:       "Test Task",
			Description: sql.NullString{String: "Test Description", Valid: true},
			Status:      domain.TaskStatusToDo,
			Priority:    domain.TaskPriorityHigh,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockQuerier := mocks.NewMockQuerier(ctrl)
		mockQuerier.EXPECT().
			GetTask(gomock.Any(), taskID).
			Return(invalidTask, nil).
			Times(1)

		service := newTestTaskService(logger, mockQuerier)

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from invalid UUID, but didn't get one")
			}
		}()

		service.GetTask(context.Background(), taskID)
	})
}