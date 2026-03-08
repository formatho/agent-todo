package cmd

import (
	"fmt"

	"github.com/formatho/agent-todo/cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Display the version, commit, build date, and builder information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("agent-todo version %s\n", config.GetVersion())
		fmt.Printf("  commit:      %s\n", config.GetCommit())
		fmt.Printf("  built at:    %s\n", config.GetDate())
		fmt.Printf("  built by:    %s\n", config.GetBuiltBy())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
