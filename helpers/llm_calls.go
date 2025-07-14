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

type LLMConfig struct {
	Backend string `yaml:"backend"`
}

type OllamaConfig struct {
	Model        string `yaml:"model"`
	Endpoint     string `yaml:"endpoint"`
	SystemPrompt string `yaml:"system_prompt"`
}

type OpenAIConfig struct {
	Model        string `yaml:"model"`
	APIKey       string `yaml:"api_key"`
	SystemPrompt string `yaml:"system_prompt"`
}

type Config struct {
	LLM    LLMConfig    `yaml:"llm"`
	Ollama OllamaConfig `yaml:"ollama"`
	OpenAI OpenAIConfig `yaml:"openai"`
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

	switch cfg.LLM.Backend {
	case "ollama":
		return RewordNoteWithOllama(&cfg.Ollama, note)
	case "openai":
		return RewordNoteWithOpenAI(&cfg.OpenAI, note)
	default:
		return "", fmt.Errorf("unsupported backend: %s", cfg.LLM.Backend)
	}

}

func RewordNoteWithOpenAI(cfg *OpenAIConfig, note string) (string, error) {
	payload := map[string]interface{}{
		"model":  cfg.Model,
		"stream": false,
		"messages": []map[string]string{
			{"role": "system", "content": cfg.SystemPrompt},
			{"role": "user", "content": fmt.Sprintf(`Original note: %q`, note)},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %s, body: %s", resp.Status, respBody)
	}

	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}
	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return res.Choices[0].Message.Content, nil
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
