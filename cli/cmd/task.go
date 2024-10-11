package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	v1 "task/pkg/gen/cloud/v1"
	cloudv1connect "task/pkg/gen/cloud/v1/cloudv1connect"
	"task/pkg/x"

	// Add this import

	"connectrpc.com/connect"
	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var taskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t"},
	Short:   "Manage tasks in the system",
	Long: `The task command allows you to manage tasks in the system, including creating new tasks, 
retrieving task details, listing all tasks, and viewing task history. 
Use subcommands to perform specific operations on tasks.`,
}

// createTaskCmd represents the create task command
var createTaskCmd = &cobra.Command{
	Use:     "create [task name] --type [task type] --parameter [key=value] --description [task description]",
	Aliases: []string{"c", "new"},
	Short:   "Create a new task",
	Long: `Create a new task in the system with the specified name, type, parameters, and description.
You must provide a task name and type. Parameters and description are optional.

The task type should be one of the predefined types in the system (e.g., send_email, run_query).
Multiple parameters can be added by repeating the --parameter flag.
The description flag allows you to add a detailed explanation of the task.`,
	Example: `  task create "Send Newsletter" --type send_email --parameter recipient=user@example.com --parameter subject="Weekly Update" --description "Send weekly newsletter to subscribers"
  task create "Generate Report" --type run_query --parameter query="SELECT * FROM sales" --parameter format=csv --description "Generate monthly sales report"
  task c "Backup Database" --type system_backup --parameter target=/backups/db.sql --description "Perform full database backup"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]
		taskType, _ := cmd.Flags().GetString("type")
		if taskType == "" {
			fmt.Println("Error: --type flag is required")
			cmd.Usage()
			os.Exit(1)
		}
		parameters, _ := cmd.Flags().GetStringToString("parameter")
		description, _ := cmd.Flags().GetString("description")
		addTask(taskName, taskType, parameters, description)
	},
}

// getTaskCmd represents the get task command
var getTaskCmd = &cobra.Command{
	Use:     "get --id [task_id]",
	Aliases: []string{"g", "show"},
	Short:   "Get details of a specific task",
	Long: `Retrieve and display the details of a specific task by its ID.
You can specify the output format as table (default), json, or yaml.`,
	Example: `  task get --id 123
  task get --id 456 --output json
  task g -i 789 -o yaml`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt64("id")
		if id == 0 {
			fmt.Println("Error: --id flag is required")
			cmd.Usage()
			os.Exit(1)
		}
		outputFormat, _ := cmd.Flags().GetString("output")
		getTask(id, outputFormat)
	},
}

// listTaskCmd represents the list task command
var listTaskCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all tasks",
	Long: `List all tasks in the system. This command displays a summary of all tasks,
including their IDs, names, types, and current statuses.
You can specify the output format as table (default), json, or yaml.
Use --offset and --limit flags for pagination, --status for filtering by status,
and --type for filtering by task type.`,
	Example: `  task list
  task list --output json
  task ls -o yaml
  task list --offset 20 --limit 10
  task list --status running
  task list --type email_send`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		outputFormat, _ := cmd.Flags().GetString("output")
		offset, _ := cmd.Flags().GetInt32("offset")
		limit, _ := cmd.Flags().GetInt32("limit")
		status, _ := cmd.Flags().GetString("status")
		taskType, _ := cmd.Flags().GetString("type")
		listTasks(outputFormat, offset, limit, status, taskType)
	},
}

// taskStatusCmd represents the task status command
var taskStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"s", "stat"},
	Short:   "Get the status counts of all tasks",
	Long:    `Retrieve and display the current status counts of all tasks in the system.`,
	Example: `  task status
  task s`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getTaskStatus()
	},
}

// init function to set up commands and flags
func init() {

	taskCmd.AddCommand(createTaskCmd, getTaskCmd, listTaskCmd, taskStatusCmd)

	addCommonFlags := func(cmd *cobra.Command) {
		cmd.Flags().Int64P("id", "i", 0, "ID of the task")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringP("output", "o", "table", "Output format (table, json, yaml)")
	}

	addCommonFlags(getTaskCmd)

	// Update flags for listTaskCmd
	listTaskCmd.Flags().StringP("output", "o", "table", "Output format (table, json, yaml)")
	listTaskCmd.Flags().Int32P("offset", "f", 0, "Offset for pagination")
	listTaskCmd.Flags().Int32P("limit", "l", 100, "Limit for pagination")
	listTaskCmd.Flags().StringP("status", "s", "all", "Filter by task status (queued, running, failed, succeeded,all)")
	listTaskCmd.Flags().StringP("type", "t", "all", "Filter by task type (e.g., email_send, run_query,all)")

	createTaskCmd.Flags().StringP("type", "t", "", "Type of the task (e.g., send_email, run_query)")
	createTaskCmd.MarkFlagRequired("type")
	createTaskCmd.Flags().StringToStringP("parameter", "p", nil, "Additional parameters for the task as key=value pairs")
	createTaskCmd.Flags().StringP("description", "d", "", "Detailed description of the task")

	rootCmd.AddCommand(taskCmd)

}

// createClient creates a new TaskManagementServiceClient with an OpenTelemetry interceptor
func createClient(address string) (cloudv1connect.TaskManagementServiceClient, error) {
	if address == "" {
		return nil, fmt.Errorf("server address is empty")
	}

	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}

	return x.CreateClient(address)
}

// addTask creates a new task and sends it to the server
func addTask(name string, taskType string, parameters map[string]string, description string) {
	slog.Info("Creating new task", "name", name, "type", taskType, "parameters", parameters, "description", description)

	client, err := createClient(address)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		return
	}
	req := connect.NewRequest(&v1.CreateTaskRequest{
		Name:        name,
		Type:        taskType,
		Description: description,
		Payload: &v1.Payload{
			Parameters: parameters,
		},
	})

	slog.Debug("Sending CreateTask request to server")
	resp, err := client.CreateTask(context.Background(), req)
	if err != nil {
		slog.Error("Error creating task", "error", err)
		return
	}

	slog.Info("Task created successfully", "id", resp.Msg.Id)
	fmt.Printf("Task created successfully:\n")
	fmt.Printf("  ID: %d\n", resp.Msg.Id)
	fmt.Printf("  Name: %s\n", name)
	fmt.Printf("  Type: %s\n", taskType)
	fmt.Printf("  Parameters: %v\n", parameters)
	fmt.Printf("  Description: %s\n", description)
}

// getTask retrieves the details of a task by its ID
func getTask(identifier int64, outputFormat string) {
	task, err := fetchTask(identifier)
	if err != nil {
		slog.Error("Error retrieving task", "error", err, "taskID", identifier)
		fmt.Printf("Error retrieving task: %v\n", err)
		return
	}
	printOutput(task, outputFormat)
}

// listTasks retrieves and displays all tasks
func listTasks(outputFormat string, offset, limit int32, status, taskType string) {
	tasks, err := fetchTasks(offset, limit, status, taskType)
	if err != nil {
		fmt.Printf("Error retrieving tasks: %v\n", err)
		return
	}
	printOutput(tasks, outputFormat)
}

// Helper function to fetch a task
func fetchTask(identifier int64) (*v1.Task, error) {

	client, err := createClient(address)
	if err != nil {

		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	req := connect.NewRequest(&v1.GetTaskRequest{Id: int32(identifier)})
	resp, err := client.GetTask(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// Helper function to fetch all tasks
func fetchTasks(offset, limit int32, status, taskType string) (*v1.TaskList, error) {
	client, err := createClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	req := &v1.TaskListRequest{
		Limit:  limit,
		Offset: offset,
	}
	// Check if status is passed and valid, if not "all" then add to request
	statusInt := x.GetStatusInt(strings.ToUpper(status))
	if statusInt != -1 {
		req.Status = &statusInt
	} else {
		return nil, fmt.Errorf("invalid status: %s", status)
	}
	// Check if type is passed and valid, if not "all" then add to request
	if taskType != "all" {
		req.Type = &taskType
	}

	resp, err := client.ListTasks(context.Background(), connect.NewRequest(req))
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

// printOutput prints the data in the specified format
func printOutput(data interface{}, format string) {
	slog.Info("Printing output", "format", format)
	switch format {
	case "table":
		x.PrintTable(data)
	case "json":
		x.PrintJSON(data)
	case "yaml":
		x.PrintYAML(data)
	default:
		slog.Warn("Invalid output format", "format", format)
		fmt.Println("Invalid output format. Use 'table', 'json', or 'yaml'.")
	}
}

var logLevel slog.Level

// InitLogger initializes the global logger with the specified log level
func InitLogger(level string) {

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		fmt.Printf("Invalid log level: %s. Using 'info' as default.\n", level)
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)
}

// getTaskStatus retrieves and displays the status counts of all tasks
func getTaskStatus() error {
	slog.Info("Retrieving task status counts")

	client, err := createClient(address)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		return fmt.Errorf("failed to create client: %w", err)
	}

	resp, err := client.GetStatus(context.Background(), connect.NewRequest(&v1.GetStatusRequest{}))
	if err != nil {
		slog.Error("Error retrieving task status counts", "error", err)
		return fmt.Errorf("error retrieving task status counts: %w", err)
	}

	fmt.Println("Task Status Counts:")
	for k, v := range resp.Msg.StatusCounts {
		statusString := x.GetStatusString(int(k))
		fmt.Printf("  %s: %d\n", statusString, v)
	}

	slog.Info("Task status counts retrieved successfully")
	return nil
}
