package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

// Top-level config struct matching your YAML layout
type Config struct {
	LLM    LLMConfig    `yaml:"llm"`
	Ollama OllamaConfig `yaml:"ollama"`
	OpenAI OpenAIConfig `yaml:"openai"`
}

// Exported sub-structs for reusability across packages
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

// Global vars used by main.go and elsewhere
var (
	ConfigPath   string  // Path to ~/.config/writeme/config.yaml
	ConfigLoaded *Config // Singleton config object loaded at startup
)

// ResolveConfigPath determines the appropriate config file path
func ResolveConfigPath() (string, error) {
	// Allow user override via WRITEME_CONFIG env var
	if configEnv := os.Getenv("WRITEME_CONFIG"); configEnv != "" {
		return filepath.Clean(configEnv), nil
	}

	var configDir string
	var err error

	if runtime.GOOS != "windows" {
		// Use XDG_CONFIG_HOME if set, else use ~/.config
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			configDir = xdg
		} else {
			var home string
			home, err = os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not get user home dir: %w", err)
			}
			configDir = filepath.Join(home, ".config")
		}
	} else {
		// Windows: use %APPDATA%, fall back to AppData/Roaming
		if appData := os.Getenv("APPDATA"); appData != "" {
			configDir = appData
		} else {
			var home string
			home, err = os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not get user home dir: %w", err)
			}
			configDir = filepath.Join(home, "AppData", "Roaming")
		}
	}

	finalPath := filepath.Join(configDir, "writeme", "config.yaml")
	return filepath.Clean(finalPath), nil
}

// LoadConfig reads and parses the config file from the given path
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
