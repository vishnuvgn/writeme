/*
Copyright © 2025 NAME HERE vishnuvgn05@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev" // default when running locally
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "writeme",
	Short: "An AI README documenter",
	Long: `An AI README editor where you can commit messages
to and have it format and append to your readme. Can connect to ai
endpoints including openai and llama.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Version: "",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// SetVersionInfo allows main.go to pass in the real version.
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d

	// This makes writeme --version show the version from GoReleaser.
	rootCmd.Version = version
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.writeme.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}
