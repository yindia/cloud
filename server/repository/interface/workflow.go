package interfaces

import (
	"context"

	model "task/server/repository/model/task"
)

// WorkflowRepo defines the interface for the task history repository.
// It handles operations related to task history management.
//
//go:generate mockery --output=../mocks --case=underscore --all --with-expecter
type WorkflowRepo interface {
	// CreateWorkflow creates a history entry for a task.
	// It takes a context.Context parameter for handling request-scoped values and deadlines.
	CreateWorkflow(ctx context.Context, workflow model.Workflow) (model.Workflow, error)

	// GetWorkflow retrieves the history of a task by its ID.
	// Returns a slice of task history entries, or an error if none found.
	GetWorkflow(ctx context.Context, workflowID uint) (*model.Workflow, error)

	// ListTaskHistories lists all history entries for a given task, with pagination support.
	// Returns a slice of task history entries, along with a pagination token (if any) for subsequent queries.
	ListWorkflow(ctx context.Context) ([]model.Workflow, error)
}
