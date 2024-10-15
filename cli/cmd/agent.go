package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	taskApi "task/controller/api/v1"
	v1 "task/pkg/gen/cloud/v1"
	"task/pkg/gen/cloud/v1/cloudv1connect"
	k8s "task/pkg/k8s"
	"task/pkg/plugins"
	"task/pkg/x"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Run the workflow orchestration server",
	Long:    `Start the workflow orchestration server and continuously stream task updates.`,
	Example: `  task serve`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWorkflowOrchestration(cmd.Context())
	},
}

// Number of worker goroutines
var numWorkers = 1000

func init() {
	rootCmd.AddCommand(serveCmd)
}

// runWorkflowOrchestration starts the workflow orchestration server and handles task updates.
func runWorkflowOrchestration(ctx context.Context) error {
	logger := slog.With("component", "workflow_orchestration")
	logger.Info("Starting workflow orchestration server")

	// Create a cancelable context for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutting down gracefully...")
			wg.Wait()
			logger.Info("Workflow orchestration server stopped")
			return nil
		default:
			if err := runStreamConnection(ctx, &wg, logger); err != nil {
				logger.Error("Stream connection error", "error", err)
				time.Sleep(5 * time.Second) // Wait before retrying
				continue
			}
		}
	}
}

func runStreamConnection(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger) error {

	var err error

	client := cloudv1connect.NewTaskManagementServiceClient(http.DefaultClient, "http://localhost:8080")
	k8sClient, err := k8s.NewK8sClient("")
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}
	go sendPeriodicRequests(ctx, logger, client) // Pass stream as a pointer

	stream, err := client.PullEvents(ctx, connect.NewRequest(&v1.PullEventsRequest{}))
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}

	for {
		ok := stream.Receive()
		if !ok {
			return fmt.Errorf("failed to receive response: %w", err)
		}

		go processWork(ctx, stream.Msg(), logger, k8sClient)
	}

	return nil
}

// sendPeriodicRequests sends periodic RunCommand requests to the server.
func sendPeriodicRequests(ctx context.Context, logger *slog.Logger, client cloudv1connect.TaskManagementServiceClient) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := client.Heartbeat(ctx, connect.NewRequest(&v1.HeartbeatRequest{
				Timestamp: time.Now().Format(time.RFC3339),
			}))

			if err != nil {
				logger.Error("Error sending request", "error", err)
				continue
			}
			logger.Debug("Sent periodic RunCommand request")
		}
	}
}

func processWork(ctx context.Context, task *v1.PullEventsResponse, logger *slog.Logger, k8sClient *k8s.K8s) {

	_, err := k8sClient.CreateTask(&taskApi.Task{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("task-%d", task.Work.Task.Id),
			Namespace: "test",
		},
		Spec: taskApi.TaskSpec{
			ID:   task.Work.Task.Id,
			Type: task.Work.Task.Type,
			Payload: taskApi.Payload{
				Parameters: task.Work.Task.Payload.Parameters,
			},
			Status:      int32(task.Work.Task.Status),
			Description: task.Work.Task.Description,
		},
	})
	if err != nil {
		logger.Error("Failed to create task", "error", err, "task", task)
		return
	}
}

// processWorkflowUpdate handles different types of responses and returns the workflow state.
func processWorkflowUpdate(ctx context.Context, task *v1.PullEventsResponse, logger *slog.Logger) (v1.TaskStatusEnum, string, error) {
	response := task.Work
	logger = logger.With("task_id", response.Task.Id)
	logger.Info("Received workflow update", "task_type", response.Task.Type)

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		logger.Info("Task processing completed", "duration", duration)
	}()

	defer func() {
		if r := recover(); r != nil {
			// Return FAILED status in case of panic
			panic(fmt.Sprintf("Task panicked: %v", r))
		}
	}()

	plugin, err := plugins.NewPlugin(response.Task.Type)
	if err != nil {
		return v1.TaskStatusEnum_FAILED, fmt.Sprintf("Failed to create plugin: %v", err), err
	}

	// Add retry logic for running the task
	runErr := plugin.Run(response.Task.Payload.Parameters)

	if runErr != nil {
		return v1.TaskStatusEnum_FAILED, fmt.Sprintf("Error running task after %d attempts: %v"), runErr
	}

	logger.Info("Task completed successfully")

	return v1.TaskStatusEnum_SUCCEEDED, "Task completed successfully", nil
}

// retry is a helper function to retry operations with exponential backoff
func retry(ctx context.Context, maxAttempts int, initialBackoff time.Duration, operation func() error) error {
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}
		if attempt == maxAttempts-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(initialBackoff * time.Duration(1<<uint(attempt))):
		}
	}
	return fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, err)
}

// updateTaskStatus updates the status of a task using the Task Management Service.
func updateTaskStatus(ctx context.Context, taskID int64, status v1.TaskStatusEnum, message string) error {
	client, err := x.CreateClient(address)
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
