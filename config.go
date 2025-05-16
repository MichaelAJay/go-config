package config

import (
	"fmt"
	"sync"
)

// Config represents the configuration interface
type Config interface {
	// Get retrieves a configuration value by key
	Get(key string) (any, bool)

	// GetString retrieves a string configuration value
	GetString(key string) (string, bool)

	// GetInt retrieves an integer configuration value
	GetInt(key string) (int, bool)

	// GetBool retrieves a boolean configuration value
	GetBool(key string) (bool, bool)

	// GetFloat retrieves a float configuration value
	GetFloat(key string) (float64, bool)

	// GetStringSlice retrieves a string slice configuration value
	GetStringSlice(key string) ([]string, bool)

	// Set sets a configuration value
	Set(key string, value any) error

	// Load loads configuration from a source
	Load(source Source) error

	// Validate validates the configuration
	Validate() error
}

// Source represents a configuration source
type Source interface {
	// Load loads configuration from the source
	Load() (map[string]any, error)
}

// ConfigManager implements the Config interface
type ConfigManager struct {
	values     map[string]any
	mu         sync.RWMutex
	validators []Validator
}

// New creates a new configuration manager
func New() *ConfigManager {
	return &ConfigManager{
		values:     make(map[string]any),
		validators: make([]Validator, 0),
	}
}

// Get retrieves a configuration value by key
func (c *ConfigManager) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.values[key]
	return value, exists
}

// GetString retrieves a string configuration value
func (c *ConfigManager) GetString(key string) (string, bool) {
	value, exists := c.Get(key)
	if !exists {
		return "", false
	}

	str, ok := value.(string)
	return str, ok
}

// GetInt retrieves an integer configuration value
func (c *ConfigManager) GetInt(key string) (int, bool) {
	value, exists := c.Get(key)
	if !exists {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetBool retrieves a boolean configuration value
func (c *ConfigManager) GetBool(key string) (bool, bool) {
	value, exists := c.Get(key)
	if !exists {
		return false, false
	}

	b, ok := value.(bool)
	return b, ok
}

// GetFloat retrieves a float configuration value
func (c *ConfigManager) GetFloat(key string) (float64, bool) {
	value, exists := c.Get(key)
	if !exists {
		return 0, false
	}

	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

// GetStringSlice retrieves a string slice configuration value
func (c *ConfigManager) GetStringSlice(key string) ([]string, bool) {
	value, exists := c.Get(key)
	if !exists {
		return nil, false
	}

	switch v := value.(type) {
	case []string:
		return v, true
	case []any:
		result := make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				result[i] = str
			} else {
				return nil, false
			}
		}
		return result, true
	default:
		return nil, false
	}
}

// Set sets a configuration value
func (c *ConfigManager) Set(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = value
	return nil
}

// Load loads configuration from a source
func (c *ConfigManager) Load(source Source) error {
	values, err := source.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range values {
		c.values[k] = v
	}

	return nil
}

// Validate validates the configuration
func (c *ConfigManager) Validate() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, validator := range c.validators {
		if err := validator.Validate(c.values); err != nil {
			return err
		}
	}

	return nil
}

// AddValidator adds a validator to the configuration
func (c *ConfigManager) AddValidator(validator Validator) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.validators = append(c.validators, validator)
}
