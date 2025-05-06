package config

import (
	"fmt"
	"reflect"
)

// Validator defines the interface for configuration validation
type Validator interface {
	Validate(values map[string]any) error
}

// RequiredValidator ensures that specified keys are present
type RequiredValidator struct {
	Keys []string
}

// Validate implements the Validator interface
func (v *RequiredValidator) Validate(values map[string]any) error {
	for _, key := range v.Keys {
		if _, exists := values[key]; !exists {
			return fmt.Errorf("required configuration key missing: %s", key)
		}
	}
	return nil
}

// TypeValidator ensures that values are of the correct type
type TypeValidator struct {
	Key  string
	Type reflect.Type
}

// Validate implements the Validator interface
func (v *TypeValidator) Validate(values map[string]any) error {
	value, exists := values[v.Key]
	if !exists {
		return nil // Skip validation if key doesn't exist
	}

	actualType := reflect.TypeOf(value)
	if actualType != v.Type {
		return fmt.Errorf("invalid type for key %s: expected %v, got %v", v.Key, v.Type, actualType)
	}
	return nil
}

// RangeValidator ensures that numeric values are within a specified range
type RangeValidator struct {
	Key   string
	Min   float64
	Max   float64
	IsInt bool
}

// Validate implements the Validator interface
func (v *RangeValidator) Validate(values map[string]any) error {
	value, exists := values[v.Key]
	if !exists {
		return nil // Skip validation if key doesn't exist
	}

	var floatValue float64
	switch val := value.(type) {
	case float64:
		floatValue = val
	case int:
		floatValue = float64(val)
	default:
		return fmt.Errorf("invalid type for range validation on key %s: expected number", v.Key)
	}

	if floatValue < v.Min || floatValue > v.Max {
		return fmt.Errorf("value for key %s is out of range: expected between %v and %v", v.Key, v.Min, v.Max)
	}

	if v.IsInt && floatValue != float64(int(floatValue)) {
		return fmt.Errorf("value for key %s must be an integer", v.Key)
	}

	return nil
}
