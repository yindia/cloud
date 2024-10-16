package interfaces

import (
	"context"

	model "task/server/repository/model/task"
)

// ExecutionRepo defines the interface for the task history repository.
// It handles operations related to task history management.
//
//go:generate mockery --output=../mocks --case=underscore --all --with-expecter
type ExecutionRepo interface {
	// CreateExecution creates a history entry for a task.
	// It takes a context.Context parameter for handling request-scoped values and deadlines.
	CreateExecution(ctx context.Context, execution model.Execution) (model.Execution, error)

	// GetExecution retrieves the history of a task by its ID.
	// Returns a slice of task history entries, or an error if none found.
	GetExecution(ctx context.Context, taskID uint) (*model.Execution, error)

	ListExecution(ctx context.Context) ([]model.Execution, error)
}
