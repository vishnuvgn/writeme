package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"writeme/config"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the writeme config file in your editor",
	Long:  `Opens the writeme config.yaml file using your default $EDITOR.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runEdit(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	configCmd.AddCommand(editCmd)
}

func runEdit() error {
	path, err := config.ResolveConfigPath()
	if err != nil {
		return fmt.Errorf("could not resolve config path: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist at %s. Run 'writeme config init' first", path)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "vi"
		}
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch editor (%s): %w", editor, err)
	}

	return nil
}
