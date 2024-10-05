package cmd

import (
	"context"
	"fmt"
	"os"

	v1 "task/pkg/gen/cloud/v1"
	"task/pkg/x"

	"connectrpc.com/connect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command for tasks
var historyCmd = &cobra.Command{
	Use:     "history --id <task_id>",
	Aliases: []string{"h", "log"},
	Short:   "Get history of a specific task",
	Long: `Retrieve and display the history of a specific task by its ID.
This command shows all status changes and events related to the task over time.
You can specify the output format as table (default), json, or yaml.`,
	Example: `  task history --id 123
  task history --id 456 --output json
  task h -i 789 -o yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt64("id")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get 'id' flag: %v\n", err)
			os.Exit(1)
		}
		if id <= 0 {
			fmt.Fprintln(os.Stderr, "Error: --id flag is required and must be a positive integer")
			cmd.Usage()
			os.Exit(1)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get 'output' flag: %v\n", err)
			os.Exit(1)
		}
		if err := getTaskHistory(id, output); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().Int64P("id", "i", 0, "ID of the task (required)")
	historyCmd.MarkFlagRequired("id")
	historyCmd.Flags().StringP("output", "o", "table", "Output format (table, json, yaml)")
}

// getTaskHistory retrieves and prints the history of a task by its ID
func getTaskHistory(identifier int64, outputFormat string) error {
	history, err := fetchTaskHistory(identifier)
	if err != nil {
		return fmt.Errorf("failed to retrieve task history: %w", err)
	}
	return printTaskHistory(history, outputFormat)
}

// fetchTaskHistory retrieves the task history from the server
func fetchTaskHistory(identifier int64) (*v1.GetTaskHistoryResponse, error) {
	client, err := x.CreateClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	req := connect.NewRequest(&v1.GetTaskHistoryRequest{Id: int32(identifier)})
	resp, err := client.GetTaskHistory(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task history: %w", err)
	}
	return resp.Msg, nil
}

// printTaskHistory prints task history in the specified format
func printTaskHistory(history *v1.GetTaskHistoryResponse, outputFormat string) error {
	switch outputFormat {
	case "json":
		return x.PrintJSON(history)
	case "yaml":
		return x.PrintYAML(history)
	case "table":
		table := tablewriter.NewWriter(os.Stdout)
		x.PrintTaskHistoryTable(table, history)
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}
