package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type OllamaConfig struct {
	Model        string `yaml:"model"`
	Endpoint     string `yaml:"endpoint"`
	SystemPrompt string `yaml:"system_prompt"`
}

type Config struct {
	Ollama OllamaConfig `yaml:"ollama"`
	// Backend string `yaml:"backend"`  // For future multi-backend support!
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal YAML: %w", err)
	}

	return &cfg, nil
}

// This is your smart wrapper.
func RewordNote(note string) (string, error) {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		return "", fmt.Errorf("could not load config.yaml: %w", err)
	}

	// Add logic for other backends here.
	return RewordNoteWithOllama(&cfg.Ollama, note)
}

// This does the actual Ollama call.
func RewordNoteWithOllama(cfg *OllamaConfig, note string) (string, error) {
	payload := map[string]interface{}{
		"model":  cfg.Model,
		"stream": false,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": cfg.SystemPrompt,
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Reword this note: \"%s\"", note),
			},
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("could not marshal payload: %w", err)
	}

	resp, err := http.Post(cfg.Endpoint, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %s, body: %s", resp.Status, string(b))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("could not decode response JSON: %w", err)
	}

	return result.Message.Content, nil
}
