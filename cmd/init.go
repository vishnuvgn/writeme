package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"writeme/config"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a blank config.yaml file",
	Long:  `Creates a blank config.yaml file in the appropriate config directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runInit()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	configCmd.AddCommand(initCmd)
}

func runInit() error {
	targetPath, err := config.ResolveConfigPath()
	if err != nil {
		return fmt.Errorf("could not resolve config path: %w", err)
	}

	if _, err := os.Stat(targetPath); err == nil {
		fmt.Printf("Config already exists at %s. Overwrite? (y/N): ", targetPath)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	defaultConfig := `llm:
  backend: ollama

ollama:
  model: llama3.1:latest
  endpoint: http://localhost:11434/api/chat
  system_prompt: |
    You are an assistant that rewrites notes for developer documentation.
    - Keep the meaning exactly the same.
    - Do NOT add new context.
    - Make it clear, concise, and direct.
    - Return exactly one line.
    - Do not say "Sure", "Here", or any greeting.

openai:
  model: gpt-4o-mini
  api_key: sk-8 # your key
  system_prompt: |
    You are an assistant that rewrites notes for developer documentation.
    Follow these rules:
    - Keep the meaning exactly the same.
    - Do not add or infer new information.
    - Make it direct and clear.
    - Output only the reworded line.
`

	if err := os.WriteFile(targetPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("could not write blank config file: %w", err)
	}

	fmt.Printf("Blank config file created at: %s\n", targetPath)
	return nil
}
