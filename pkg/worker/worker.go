package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	v1 "task/pkg/gen/cloud/v1"
	"task/pkg/gen/cloud/v1/cloudv1connect"
	"task/pkg/plugins"
	"task/pkg/x"
	"task/server/repository/model/task"
	"time"

	"connectrpc.com/connect"
	"github.com/riverqueue/river"
)

var cloudClient cloudv1connect.TaskManagementServiceClient

// TaskArgument represents the argument structure for a task job.
type TaskArgument struct {
	Task task.Task `json:"task"`
}

// Kind returns the kind of the task argument.
func (TaskArgument) Kind() string { return "email_send" }

// InsertOpts returns the insertion options for the task.
func (TaskArgument) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 5}
}

// TaskWorker is the worker implementation for processing tasks.
type TaskWorker struct {
	river.WorkerDefaults[TaskArgument]
}

// Work processes a single task job.
func (w *TaskWorker) Work(ctx context.Context, job *river.Job[TaskArgument]) error {
	logger := slog.With("task_id", job.Args.Task.ID, "attempt", job.Attempt)
	logger.Info("Starting task processing")

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		logger.Info("Task processing completed", "duration", duration, "task_type", job.Args.Task.Type)
	}()

	if err := updateTaskStatus(ctx, int64(job.Args.Task.ID), v1.TaskStatusEnum_RUNNING, fmt.Sprintf("Task started (Attempt %d)", job.Attempt)); err != nil {
		logger.Error("Failed to update task status to RUNNING", "error", err)
		return fmt.Errorf("failed to update task status to RUNNING: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Task panicked", "panic", r)
			if err := updateTaskStatus(ctx, int64(job.Args.Task.ID), v1.TaskStatusEnum_FAILED, fmt.Sprintf("Task panicked (Attempt %d): %v", job.Attempt, r)); err != nil {
				logger.Error("Failed to update task status after panic", "error", err)
			}
		}
	}()

	plugin, err := plugins.NewPlugin(job.Args.Task.Type)
	if err != nil {
		return w.handleError(ctx, job, logger, "Failed to create plugin", err)
	}

	var payloadMap map[string]string
	if err := json.Unmarshal([]byte(job.Args.Task.Payload), &payloadMap); err != nil {
		return w.handleError(ctx, job, logger, "Failed to unmarshal payload", err)
	}

	if err := plugin.Run(payloadMap); err != nil {
		return w.handleError(ctx, job, logger, "Error running task", err)
	}

	logger.Info("Task completed successfully")
	if err := updateTaskStatus(ctx, int64(job.Args.Task.ID), v1.TaskStatusEnum_SUCCEEDED, fmt.Sprintf("Task completed successfully (Attempt %d)", job.Attempt)); err != nil {
		logger.Error("Failed to update task status to SUCCEEDED", "error", err)
		return fmt.Errorf("failed to update task status to SUCCEEDED: %w", err)
	}

	return nil
}

// handleError is a helper function to handle errors during task processing.
func (w *TaskWorker) handleError(ctx context.Context, job *river.Job[TaskArgument], logger *slog.Logger, message string, err error) error {
	logger.Error(message, "error", err)
	errorMsg := fmt.Sprintf("%s (Attempt %d): %v", message, job.Attempt, err)
	if updateErr := updateTaskStatus(ctx, int64(job.Args.Task.ID), v1.TaskStatusEnum_FAILED, errorMsg); updateErr != nil {
		logger.Error("Failed to update task status to FAILED", "error", updateErr)
	}
	return fmt.Errorf("%s for task %d (Attempt %d): %w", message, job.Args.Task.ID, job.Attempt, err)
}

// NextRetry determines the time for the next retry attempt.
func (w *TaskWorker) NextRetry(job *river.Job[TaskArgument]) time.Time {
	return time.Now().Add(2 * time.Second)
}

// Timeout sets the maximum duration for a task to complete.
func (w *TaskWorker) Timeout(job *river.Job[TaskArgument]) time.Duration {
	timeout := 10
	if timeoutStr := os.Getenv("TASK_TIME_OUT"); timeoutStr != "" {
		if parsedTimeout, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	return (time.Duration(timeout) + 5) * time.Second
}

// updateTaskStatus updates the status of a task using the Task Management Service.
func updateTaskStatus(ctx context.Context, taskID int64, status v1.TaskStatusEnum, message string) error {
	client, err := x.CreateClient(os.Getenv("SERVER_ENDPOINT"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	_, err = client.UpdateTaskStatus(ctx, connect.NewRequest(&v1.UpdateTaskStatusRequest{
		Id:      int32(taskID),
		Status:  status,
		Message: message,
	}))
	if err != nil {
		return fmt.Errorf("failed to update task %d status: %w", taskID, err)
	}

	return nil
}
