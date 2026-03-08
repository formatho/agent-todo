package main

import (
	"fmt"
	"os"

	"github.com/formatho/agent-todo/cli/cmd"
	"github.com/formatho/agent-todo/cli/config"
)

var (
	// Version information (set by build flags)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

func main() {
	// Set version in config
	config.SetVersion(Version)
	config.SetCommit(Commit)
	config.SetDate(Date)
	config.SetBuiltBy(BuiltBy)

	cmd.Execute()
}

func init() {
	// Add version flag
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("agent-todo version %s\n", Version)
		fmt.Printf("commit: %s\n", Commit)
		fmt.Printf("built at: %s by %s\n", Date, BuiltBy)
		os.Exit(0)
	}
}
