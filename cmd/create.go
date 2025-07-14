package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a NOTES.md file with your project folder name as the heading",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := "NOTES.md"

		// Check if file already exists
		if _, err := os.Stat(fileName); err == nil {
			// File exists — ask user
			fmt.Printf("%s already exists. Overwrite? (y/n): ", fileName)
			reader := bufio.NewReader(os.Stdin)
			answer, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("could not read input: %w", err)
			}
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Canceled.")
				return nil
			}
		}

		// Get the current directory name
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get working directory: %w", err)
		}
		dirName := filepath.Base(wd)

		// Create (or overwrite) the file
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		defer f.Close()

		// Write the heading with newline
		_, err = fmt.Fprintf(f, "# %s\n", dirName)
		if err != nil {
			return fmt.Errorf("could not write to file: %w", err)
		}

		fmt.Printf("%s created with heading '# %s'\n", fileName, dirName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
