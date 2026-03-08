package cmd

import (
	"fmt"
	"os"

	"github.com/formatho/agent-todo/cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent-todo",
	Short: "CLI for Agent Todo Management Platform",
	Long: `A comprehensive CLI tool for managing projects, tasks, and AI agents
in the Agent Todo Management Platform.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for completion command
		if cmd.Name() == "completion" || cmd.Name() == "help" {
			return nil
		}

		// Initialize config
		if err := config.Init(); err != nil {
			return fmt.Errorf("error initializing config: %w", err)
		}

		// Apply server URL flag if provided
		if serverURL != "" {
			config.SetServerURL(serverURL)
		}

		return nil
	},
}

var (
	serverURL string
	verbose   bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "", "Server URL (overrides config file)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
