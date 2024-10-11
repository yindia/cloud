package interfaces

import (
	"context"

	model "task/server/repository/model/task"
)

// TaskRepo defines the interface for the task repository.
// It handles operations related to task management, including task creation, status update, and history retrieval.
//
//go:generate mockery --output=../mocks --case=underscore --all --with-expecter
type TaskRepo interface {
	// CreateTask creates a new task with the provided information.
	// It takes a context.Context parameter for handling request-scoped values and deadlines.
	CreateTask(ctx context.Context, task model.Task) (model.Task, error)

	// GetTaskByID retrieves a task by its ID.
	// It returns the task if found, or an error otherwise.
	GetTaskByID(ctx context.Context, taskID uint) (*model.Task, error)

	// UpdateTaskStatus updates the status of a task.
	// It requires the task ID and the new status to be set.
	UpdateTaskStatus(ctx context.Context, taskID uint, status int) error

	// ListTasks retrieves a list of tasks based on the provided criteria.
	// It takes a context.Context parameter for handling request-scoped values and deadlines.
	// The limit and offset parameters are used for pagination.
	// The status parameter filters tasks by their status (use -1 for all statuses).
	// The taskType parameter filters tasks by their type (use an empty string for all types).
	// It returns a slice of tasks and an error if any occurs during the operation.
	ListTasks(ctx context.Context, limit, offset int, status int, taskType string) ([]model.Task, error)

	// GetTaskStatusCounts retrieves the count of tasks for each status.
	// It takes a context.Context parameter for handling request-scoped values and deadlines.
	// It returns a map where the key is the status code and the value is the count of tasks with that status.
	// An error is returned if any occurs during the operation.
	GetTaskStatusCounts(ctx context.Context) (map[int]int64, error)

	GetStalledTasks(ctx context.Context) ([]model.Task, error)
}
