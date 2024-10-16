package cmd

import (
	"context"
	"fmt"
	cloudv1 "task/pkg/gen/cloud/v1"
	v1 "task/pkg/gen/cloud/v1"
	"task/pkg/gen/cloud/v1/cloudv1connect"
	"task/pkg/x"
	"time"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
)

var (
	end2endCmd = &cobra.Command{
		Use:   "end2end",
		Short: "Run end-to-end tests for the system",
		Long: `This command executes a series of end-to-end tests to verify the entire system's functionality.
It creates a specified number of tasks and monitors their completion status for up to 3 minutes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnd2EndTests()
		},
	}
	numTasks int
)

func init() {
	rootCmd.AddCommand(end2endCmd)
	end2endCmd.Flags().IntVarP(&numTasks, "count", "n", 300, "Number of tasks to create (default 100, max 100)")
}

func runEnd2EndTests() error {
	fmt.Println("Starting end-to-end tests...")
	client, err := x.CreateClient(address)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	initialStatus, err := checkInitialStatus(client)
	if err != nil {
		return fmt.Errorf("failed to check initial status: %w", err)
	}

	if err := createTasks(client); err != nil {
		return fmt.Errorf("failed to create tasks: %w", err)
	}

	return monitorTasks(client, initialStatus)
}

func checkInitialStatus(client cloudv1connect.TaskManagementServiceClient) (map[v1.TaskStatusEnum]int64, error) {
	fmt.Println("Checking initial task status...")
	response, err := checkTaskStatus(client)
	if err != nil {
		return nil, err
	}

	initialStatus := make(map[v1.TaskStatusEnum]int64)
	for status, count := range response.Msg.StatusCounts {
		initialStatus[v1.TaskStatusEnum(status)] = count
	}

	fmt.Printf("Initial status: %d queued, %d running, %d succeeded, %d failed\n",
		initialStatus[v1.TaskStatusEnum_QUEUED],
		initialStatus[v1.TaskStatusEnum_RUNNING],
		initialStatus[v1.TaskStatusEnum_SUCCEEDED],
		initialStatus[v1.TaskStatusEnum_FAILED])
	return initialStatus, nil
}

func createTasks(client cloudv1connect.TaskManagementServiceClient) error {
	fmt.Printf("Creating %d tasks in parallel...\n", numTasks)

	// Create a buffered channel to limit concurrency to 100
	semaphore := make(chan struct{}, 100)
	errChan := make(chan error, numTasks)

	for i := 0; i < numTasks; i++ {
		semaphore <- struct{}{} // Acquire semaphore
		go func(index int) {
			defer func() { <-semaphore }() // Release semaphore

			taskType := getTaskType(index)
			if err := createSingleTask(client, index, taskType); err != nil {
				errChan <- fmt.Errorf("failed to create task %d: %w", index+1, err)
			} else {
				errChan <- nil
			}
		}(i)
	}

	// Wait for all goroutines to finish and check for errors
	createdTasks := 0
	for i := 0; i < numTasks; i++ {
		if err := <-errChan; err != nil {
			return err
		}
		createdTasks++
		if (createdTasks%50 == 0) || (createdTasks == numTasks) {
			fmt.Printf("Progress: Created %d/%d tasks\n", createdTasks, numTasks)
		}
	}

	fmt.Printf("Successfully created %d tasks. Now monitoring...\n", createdTasks)
	return nil
}

func getTaskType(index int) string {
	if index < numTasks/2 { // 50% of tasks will be run_query
		return "run_query"
	}
	return "send_email" // The remaining 50% will be send_email
}

func createSingleTask(client cloudv1connect.TaskManagementServiceClient, index int, taskType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.CreateTask(ctx, connect.NewRequest(&v1.CreateTaskRequest{
		Name:        fmt.Sprintf("Task %d", index+1),
		Description: fmt.Sprintf("Description for Task %d", index+1),
		Type:        taskType,
		Payload: &v1.Payload{
			Parameters: map[string]string{
				"test": fmt.Sprintf("test_%d", index+1),
			},
		},
	}))

	if err != nil {
		return fmt.Errorf("failed to create task %d: %w", index+1, err)
	}
	return nil
}

func monitorTasks(client cloudv1connect.TaskManagementServiceClient, initialStatus map[v1.TaskStatusEnum]int64) error {

	fmt.Println("Monitoring task completion...")
	startTime := time.Now()
	duration := 3 * time.Minute
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	initialTotal := int64(0)
	for _, count := range initialStatus {
		initialTotal += count
	}

	for {
		response, err := checkTaskStatus(client)
		if err != nil {
			return err
		}

		currentTotal := int64(0)
		completed := int64(0)
		failed := int64(0)
		inProgress := int64(0)

		for status, count := range response.Msg.StatusCounts {
			currentTotal += count
			switch v1.TaskStatusEnum(status) {
			case v1.TaskStatusEnum_QUEUED, v1.TaskStatusEnum_RUNNING:
				inProgress += count
			case v1.TaskStatusEnum_SUCCEEDED:
				completed += count
			case v1.TaskStatusEnum_FAILED:
				failed += count
			}
		}

		newTasks := currentTotal - initialTotal

		fmt.Printf("Progress: %d new tasks, %d succeeded, %d failed, %d in progress\n",
			newTasks, completed, failed, inProgress)

		if completed+failed == int64(numTasks) {
			fmt.Printf("All tasks finished in %s. %d succeeded, %d failed.\n",
				time.Since(startTime).Round(time.Second), completed, failed)
			return nil
		}

		if time.Since(startTime) > duration {
			fmt.Printf("Test completed after %s. %d succeeded, %d failed, %d still in progress.\n",
				duration, completed, failed, inProgress)
			return nil
		}

		<-ticker.C // Wait for the next tick
	}
}
func checkTaskStatus(client cloudv1connect.TaskManagementServiceClient) (*connect.Response[cloudv1.GetStatusResponse], error) {
	resp, err := client.GetStatus(context.Background(), connect.NewRequest(&v1.GetStatusRequest{}))
	if err != nil {
		return nil, fmt.Errorf("error retrieving task status counts: %w", err)
	}

	return resp, nil
}

// Add this helper function at the end of the file
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
