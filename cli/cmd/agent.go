package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	v1 "task/pkg/gen/cloud/v1"
	"task/pkg/plugins"
	"task/pkg/x"
	"time"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	client := v1.NewTaskManagementServiceClient(conn)

	// Create a buffered channel for work assignments
	workChan := make(chan *v1.StreamResponse_WorkAssignment, 100)

	// Start the stream
	stream, err := client.StreamConnection(ctx)
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}

	// Start goroutines
	wg.Add(3)
	go sendPeriodicRequests(ctx, stream, logger, wg)
	go receiveResponses(ctx, stream, workChan, logger, wg)
	go processWorkAssignments(ctx, stream, workChan, logger, wg)

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}

// sendPeriodicRequests sends periodic RunCommand requests to the server.
func sendPeriodicRequests(ctx context.Context, stream v1.TaskManagementService_StreamConnectionClient, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := stream.Send(&v1.StreamRequest{Request: &v1.StreamRequest_Heartbeat{Heartbeat: &v1.Heartbeat{}}}); err != nil {
				logger.Error("Error sending request", "error", err)
				return
			}
			logger.Debug("Sent periodic RunCommand request")
		}
	}
}

func receiveResponses(ctx context.Context, stream v1.TaskManagementService_StreamConnectionClient, workChan chan<- *v1.StreamResponse_WorkAssignment, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF || ctx.Err() != nil {
					logger.Info("Stream closed")
					return
				}
				logger.Error("Error receiving response", "error", err)
				time.Sleep(time.Second) // Wait before retrying
				continue
			}
			fmt.Println(response)
			switch resp := response.Response.(type) {
			case *v1.StreamResponse_WorkAssignment:
				select {
				case workChan <- resp:
					logger.Debug("Work assignment queued", "task_id", resp.WorkAssignment.Task.Id)
				default:
					logger.Warn("Work channel full, discarding work assignment", "task_id", resp.WorkAssignment.Task.Id)
				}
			default:
				logger.Warn("Received unknown response type", "type", fmt.Sprintf("%T", resp))
			}
		}
	}
}

func processWorkAssignments(ctx context.Context, stream v1.TaskManagementService_StreamConnectionClient, workChan <-chan *v1.StreamResponse_WorkAssignment, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case work := <-workChan:
			status, message, err := processWorkflowUpdate(ctx, work, logger)
			if err != nil {
				logger.Error("Error processing workflow update", "error", err)
			}
			// Send status update through the stream with retries
			if err := sendStatusUpdateWithRetry(ctx, stream, work.WorkAssignment.Task.Id, status, message); err != nil {
				logger.Error("Failed to send status update after retries", "error", err)
			}
		}
	}
}

func sendStatusUpdateWithRetry(ctx context.Context, stream v1.TaskManagementService_StreamConnectionClient, taskID int32, status v1.TaskStatusEnum, message string) error {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := stream.Send(&v1.StreamRequest{
			Request: &v1.StreamRequest_UpdateTaskStatus{
				UpdateTaskStatus: &v1.UpdateTaskStatusRequest{
					Id:      taskID,
					Status:  status,
					Message: message,
				},
			},
		})
		if err == nil {
			return nil
		}
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return fmt.Errorf("failed to send status update after %d retries", maxRetries)
}

// processWorkflowUpdate handles different types of responses and returns the workflow state.
func processWorkflowUpdate(ctx context.Context, work *v1.StreamResponse_WorkAssignment, logger *slog.Logger) (v1.TaskStatusEnum, string, error) {
	response := work.WorkAssignment
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

	if err := plugin.Run(response.Task.Payload.Parameters); err != nil {
		return v1.TaskStatusEnum_FAILED, fmt.Sprintf("Error running task: %v", err), err
	}

	logger.Info("Task completed successfully")
	if err := updateTaskStatus(ctx, int64(response.Task.Id), v1.TaskStatusEnum_SUCCEEDED, "Task completed successfully"); err != nil {
		logger.Error("Failed to update task status to SUCCEEDED", "error", err)
		return v1.TaskStatusEnum_FAILED, fmt.Sprintf("Failed to update task status to SUCCEEDED: %v", err), err
	}

	return v1.TaskStatusEnum_SUCCEEDED, "", nil
}

// handlePanic recovers from panics and updates the task status accordingly.
func handlePanic(ctx context.Context, taskID int32, logger *slog.Logger) {
	if r := recover(); r != nil {
		logger.Error("Task panicked", "panic", r)
		if err := updateTaskStatus(ctx, int64(taskID), v1.TaskStatusEnum_FAILED, fmt.Sprintf("Task panicked: %v", r)); err != nil {
			logger.Error("Failed to update task status after panic", "error", err)
		}
	}
}

// handleError updates the task status to FAILED and logs the error.
func handleError(ctx context.Context, taskID int32, logger *slog.Logger, message string, err error) error {
	logger.Error(message, "error", err)
	if updateErr := updateTaskStatus(ctx, int64(taskID), v1.TaskStatusEnum_FAILED, fmt.Sprintf("%s: %v", message, err)); updateErr != nil {
		logger.Error("Failed to update task status to FAILED", "error", updateErr)
	}
	return fmt.Errorf("%s: %w", message, err)
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
