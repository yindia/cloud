package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var (
	LogLevel string
	address  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task-cli",
	Short: "A CLI application for managing task in a distributed system",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, display the help message
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you can define flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.task-cli.yaml)")

	// Add global flag for log level
	rootCmd.PersistentFlags().StringVar(&LogLevel, "log-level", "error", "Set the logging level")

	// Add global flag for address
	rootCmd.PersistentFlags().StringVar(&address, "address", "http://127.0.0.1:8080", "Set the server address")

	// Modify the PersistentPreRun function
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		InitLogger(LogLevel)
	}
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
