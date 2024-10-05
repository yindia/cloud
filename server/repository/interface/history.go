package interfaces

import (
	"context"

	model "task/server/repository/model/task"
)

// TaskHistoryRepo defines the interface for the task history repository.
// It handles operations related to task history management.
//
//go:generate mockery --output=../mocks --case=underscore --all --with-expecter
type TaskHistoryRepo interface {
	// CreateTaskHistory creates a history entry for a task.
	// It takes a context.Context parameter for handling request-scoped values and deadlines.
	CreateTaskHistory(ctx context.Context, history model.TaskHistory) (model.TaskHistory, error)

	// GetTaskHistory retrieves the history of a task by its ID.
	// Returns a slice of task history entries, or an error if none found.
	GetTaskHistory(ctx context.Context, taskID uint) ([]model.TaskHistory, error)

	// ListTaskHistories lists all history entries for a given task, with pagination support.
	// Returns a slice of task history entries, along with a pagination token (if any) for subsequent queries.
	ListTaskHistories(ctx context.Context, taskID uint) ([]model.TaskHistory, error)
}
