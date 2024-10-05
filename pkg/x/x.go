package x

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"task/pkg/config"
	cloudv1 "task/pkg/gen/cloud/v1"
	"task/pkg/gen/cloud/v1/cloudv1connect"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

var CRON_TIME = 30

// loadEnv loads environment variables from .env file
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		// Instead of returning an error, we'll just log a message
		fmt.Println("No .env file found, proceeding with default values")
	}
	return nil
}

// loadConfig processes environment variables into a Config struct
func LoadConfig() (config.Config, error) {
	var env config.Config
	if err := envconfig.Process("", &env); err != nil {
		return env, fmt.Errorf("error loading environment variables: %w", err)
	}
	return env, nil
}

// Example function to convert a map[string]string to a JSON string
func ConvertMapToJson(parameters map[string]string) (string, error) {
	// Marshal the map into a JSON byte array
	jsonBytes, err := json.Marshal(parameters)
	if err != nil {
		return "", fmt.Errorf("error converting map to JSON: %v", err)
	}

	// Convert the byte array to a string and return
	return string(jsonBytes), nil
}

// ConvertJsonToMap converts a JSON string into a map[string]string.
func ConvertJsonToMap(jsonString string) (map[string]string, error) {
	var result map[string]string

	// Unmarshal the JSON string into the map
	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return nil, fmt.Errorf("error converting JSON to map: %v", err)
	}

	return result, nil
}

// GetStatusString converts a status number to its corresponding string representation.
func GetStatusString(status int) string {
	switch status {
	case 0:
		return "QUEUED"
	case 1:
		return "RUNNING"
	case 2:
		return "FAILED"
	case 3:
		return "SUCCEEDED"
	default:
		return "UNKNOWN"
	}
}

// GetStatusInt converts a status string to its corresponding integer value.
func GetStatusInt(status string) cloudv1.TaskStatusEnum {
	switch status {
	case "QUEUED":
		return cloudv1.TaskStatusEnum_QUEUED
	case "RUNNING":
		return cloudv1.TaskStatusEnum_RUNNING
	case "FAILED":
		return cloudv1.TaskStatusEnum_FAILED
	case "SUCCEEDED":
		return cloudv1.TaskStatusEnum_SUCCEEDED
	default:
		return cloudv1.TaskStatusEnum_ALL // Indicating an unknown status
	}
}

// CreateClient creates a new TaskManagementServiceClient with an OpenTelemetry interceptor
func CreateClient(address string) (cloudv1connect.TaskManagementServiceClient, error) {
	interceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, fmt.Errorf("error creating interceptor: %w", err)
	}

	slog.Info("Creating TaskManagementServiceClient", "serverURL", address)
	return cloudv1connect.NewTaskManagementServiceClient(
		http.DefaultClient,
		address,
		connect.WithInterceptors(interceptor),
	), nil
}

// printTaskTable prints a single task in a table format
func PrintTaskTable(table *tablewriter.Table, task *cloudv1.Task) {
	table.SetHeader([]string{"Field", "Value"})
	table.Append([]string{"ID", fmt.Sprintf("%d", task.Id)})
	table.Append([]string{"Name", task.Name})
	table.Append([]string{"Type", task.Type})
	table.Append([]string{"Status", task.Status.String()})
	table.Append([]string{"Description", task.Description})
	table.Render()
}

// printTaskListTable prints a list of tasks in a table format
func PrintTaskListTable(table *tablewriter.Table, tasks *cloudv1.TaskList) {
	table.SetHeader([]string{"ID", "Name", "Type", "Status", "Description"})
	for _, task := range tasks.Tasks {
		table.Append([]string{
			fmt.Sprintf("%d", task.Id),
			task.Name,
			task.Type,
			task.Status.String(),
			truncateMessage(task.Description),
		})
	}
	table.Render()
}

// printJSON prints data in JSON format
func PrintJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

// printYAML prints data in YAML format
func PrintYAML(data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling to YAML: %v\n", err)
		return err
	}
	fmt.Println(string(yamlData))
	return nil
}

// printTaskHistoryTable prints task history in a table format
func PrintTaskHistoryTable(table *tablewriter.Table, history *cloudv1.GetTaskHistoryResponse) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	table.SetHeader([]string{"ID", "Created At", "Message", "Status"})
	for _, entry := range history.History {
		createdAt, _ := time.Parse(time.RFC3339, entry.CreatedAt)
		table.Append([]string{
			fmt.Sprintf("%d", entry.Id),
			createdAt.Format(time.RFC3339),
			truncateMessage(entry.Details),
			entry.Status.String(),
		})
	}
	table.Render()
}

// Helper function to truncate message to 30 characters
func truncateMessage(message string) string {
	if len(message) > 30 {
		return message[:27] + "..."
	}
	return message
}

// printTable prints data in table format using github.com/olekukonko/tablewriter
func PrintTable(data interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	switch v := data.(type) {
	case *cloudv1.Task:
		PrintTaskTable(table, v)
	case *cloudv1.TaskList:
		PrintTaskListTable(table, v)
	default:
		log.Println("Unsupported data type for table format")
		fmt.Println("Unsupported data type for table format")
	}
}
