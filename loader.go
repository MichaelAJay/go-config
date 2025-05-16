package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileSource loads configuration from a file
type FileSource struct {
	Path string
}

// Load implements the Source interface
func (s *FileSource) Load() (map[string]any, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(s.Path))
	var values map[string]any

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &values); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &values); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	return values, nil
}

// EnvSource loads configuration from environment variables
type EnvSource struct {
	Prefix string
}

// Load implements the Source interface
func (s *EnvSource) Load() (map[string]any, error) {
	values := make(map[string]any)
	prefix := strings.ToUpper(s.Prefix)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// Convert the key to lowercase and remove the prefix
		configKey := strings.ToLower(strings.TrimPrefix(key, prefix))
		// Replace underscores with dots for hierarchical config
		configKey = strings.ReplaceAll(configKey, "_", ".")

		// Try to parse the value as different types
		if value == "true" || value == "false" {
			values[configKey] = value == "true"
		} else if intVal, err := parseInt(value); err == nil {
			values[configKey] = intVal
		} else if floatVal, err := parseFloat(value); err == nil {
			values[configKey] = floatVal
		} else {
			values[configKey] = value
		}
	}

	return values, nil
}

// DefaultSource provides default configuration values
type DefaultSource struct {
	Values map[string]any
}

// Load implements the Source interface
func (s *DefaultSource) Load() (map[string]any, error) {
	return s.Values, nil
}

// Helper functions for parsing values
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
