package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds application configuration.
type Config struct {
	Port          int    `yaml:"port"`
	GastownDir    string `yaml:"gastown-dir"`
	MulticlaudeDir string `yaml:"multiclaude-dir"`
}

// DefaultPath returns the default config file path.
func DefaultPath() string {
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, "agent-chat", "config.yaml")
	}
	// Fallback for systems where UserConfigDir fails
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "agent-chat", "config.yaml")
	}
	return ""
}

// Load reads configuration from a YAML file.
// If path is empty, it uses the default path.
// If the file doesn't exist and path was not explicitly specified, returns empty config with no error.
// If the file doesn't exist and path was explicitly specified, returns an error.
func Load(path string, explicit bool) (*Config, error) {
	if path == "" {
		path = DefaultPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if explicit {
				return nil, fmt.Errorf("config file not found: %s", path)
			}
			// Default path doesn't exist - that's fine
			return &Config{}, nil
		}
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Errorf("permission denied reading config file: %s\nCheck file permissions with: ls -la %s", path, path)
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		var yamlErr *yaml.TypeError
		if errors.As(err, &yamlErr) {
			return nil, fmt.Errorf("invalid config file %s: %s", path, yamlErr.Errors[0])
		}
		return nil, fmt.Errorf("invalid YAML in config file %s: %w", path, err)
	}

	return &cfg, nil
}
