package config_test

import (
	"errors"
	"testing"

	"github.com/MichaelAJay/go-config"
)

// MockSource implements the Source interface for testing
type MockSource struct {
	values map[string]any
	err    error
}

func (m *MockSource) Load() (map[string]any, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.values, nil
}

// MockValidator implements the Validator interface for testing
type MockValidator struct {
	shouldFail bool
}

func (m *MockValidator) Validate(values map[string]any) error {
	if m.shouldFail {
		return errors.New("validation failed")
	}
	return nil
}

func TestConfigManager_BasicOperations(t *testing.T) {
	cfg := config.New()

	// Test setting and getting string
	err := cfg.Set("string_key", "test_value")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetString("string_key"); !ok || val != "test_value" {
		t.Errorf("GetString failed: got %v, %v, want %v, %v", val, ok, "test_value", true)
	}

	// Test setting and getting int
	err = cfg.Set("int_key", 42)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetInt("int_key"); !ok || val != 42 {
		t.Errorf("GetInt failed: got %v, %v, want %v, %v", val, ok, 42, true)
	}

	// Test setting and getting bool
	err = cfg.Set("bool_key", true)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetBool("bool_key"); !ok || val != true {
		t.Errorf("GetBool failed: got %v, %v, want %v, %v", val, ok, true, true)
	}

	// Test setting and getting float
	err = cfg.Set("float_key", 3.14)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetFloat("float_key"); !ok || val != 3.14 {
		t.Errorf("GetFloat failed: got %v, %v, want %v, %v", val, ok, 3.14, true)
	}

	// Test setting and getting string slice
	err = cfg.Set("slice_key", []string{"one", "two", "three"})
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetStringSlice("slice_key"); !ok || len(val) != 3 {
		t.Errorf("GetStringSlice failed: got %v, %v, want %v, %v", val, ok, []string{"one", "two", "three"}, true)
	}
}

func TestConfigManager_TypeConversion(t *testing.T) {
	cfg := config.New()

	// Test int to float conversion
	err := cfg.Set("number", 42)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetFloat("number"); !ok || val != 42.0 {
		t.Errorf("GetFloat failed for int: got %v, %v, want %v, %v", val, ok, 42.0, true)
	}

	// Test float to int conversion
	err = cfg.Set("float", 42.7)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	if val, ok := cfg.GetInt("float"); !ok || val != 42 {
		t.Errorf("GetInt failed for float: got %v, %v, want %v, %v", val, ok, 42, true)
	}
}

func TestConfigManager_Load(t *testing.T) {
	cfg := config.New()

	// Test successful load
	source := &MockSource{
		values: map[string]any{
			"key1": "value1",
			"key2": 42,
		},
	}
	err := cfg.Load(source)
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}
	if val, ok := cfg.GetString("key1"); !ok || val != "value1" {
		t.Errorf("GetString failed after load: got %v, %v, want %v, %v", val, ok, "value1", true)
	}

	// Test load error
	errorSource := &MockSource{
		err: errors.New("load error"),
	}
	err = cfg.Load(errorSource)
	if err == nil {
		t.Error("Expected error from Load, got nil")
	}
}

func TestConfigManager_Validation(t *testing.T) {
	cfg := config.New()

	// Test successful validation
	validator := &MockValidator{shouldFail: false}
	cfg.AddValidator(validator)
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}

	// Test failed validation
	failingValidator := &MockValidator{shouldFail: true}
	cfg.AddValidator(failingValidator)
	err = cfg.Validate()
	if err == nil {
		t.Error("Expected validation error, got nil")
	}
}

func TestConfigManager_NonExistentKeys(t *testing.T) {
	cfg := config.New()

	// Test getting non-existent keys
	if _, ok := cfg.GetString("nonexistent"); ok {
		t.Error("GetString should return false for non-existent key")
	}
	if _, ok := cfg.GetInt("nonexistent"); ok {
		t.Error("GetInt should return false for non-existent key")
	}
	if _, ok := cfg.GetBool("nonexistent"); ok {
		t.Error("GetBool should return false for non-existent key")
	}
	if _, ok := cfg.GetFloat("nonexistent"); ok {
		t.Error("GetFloat should return false for non-existent key")
	}
	if _, ok := cfg.GetStringSlice("nonexistent"); ok {
		t.Error("GetStringSlice should return false for non-existent key")
	}
}
