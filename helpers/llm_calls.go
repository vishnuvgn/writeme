package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"writeme/config"
)

// This is your smart wrapper.
func RewordNote(cfg *config.Config, note string) (string, error) {
	switch cfg.LLM.Backend {
	case "ollama":
		return RewordNoteWithOllama(&cfg.Ollama, note)
	case "openai":
		return RewordNoteWithOpenAI(&cfg.OpenAI, note)
	default:
		return "", fmt.Errorf("unsupported backend: %s", cfg.LLM.Backend)
	}
}

func RewordNoteWithOpenAI(cfg *config.OpenAIConfig, note string) (string, error) {
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
func RewordNoteWithOllama(cfg *config.OllamaConfig, note string) (string, error) {
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
