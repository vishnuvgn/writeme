package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"  // default for local builds
	commit  = "none" // default
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "writeme",
	Short:   "An AI README documenter",
	Long:    `An AI README editor where you can commit messages to and have it format and append to your readme. Can connect to AI endpoints including OpenAI and Llama.`,
	Version: version, // Use the package var directly
}

// Execute runs the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Nothing needed here. Cobra handles --version automatically.
}
