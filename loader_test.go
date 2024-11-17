package sleego

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoader_Load_Success(t *testing.T) {
	// Create a temporary JSON file with valid content
	content := `[
		{
			"name": "TestApp",
			"allowed_from": "08:00",
			"allowed_to": "18:00"
		},
		{
			"name": "TestApp2",
			"allowed_from": "09:00",
			"allowed_to": "20:00"
		}
	]`
	tmpfile, err := os.CreateTemp("./", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	loader := &Loader{}
	config, err := loader.Load(tmpfile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	app1 := config[0]
	app2 := config[1]
	for i, app := range []AppConfig{app1, app2} {
		expectedName := "TestApp"
		expectedFrom := "08:00"
		expectedTo := "18:00"
		if i == 1 {
			expectedName = "TestApp2"
			expectedFrom = "09:00"
			expectedTo = "20:00"
		}
		if app.Name != expectedName {
			t.Errorf("Expected Name to be '%s', got '%s'", expectedName, app.Name)
		}
		if app.AllowedFrom != expectedFrom {
			t.Errorf("Expected AllowedFrom to be '%s', got '%s'", expectedFrom, app.AllowedFrom)
		}
		if app.AllowedTo != expectedTo {
			t.Errorf("Expected AllowedTo to be '%s', got '%s'", expectedTo, app.AllowedTo)
		}
	}
}

func TestLoader_Load_FileNotFound(t *testing.T) {
	loader := &Loader{}
	_, err := loader.Load("nonexistent.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoader_Load_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	content := `{"name": "TestApp", "allowed_from": "08:00", "allowed_to": "18:00"` // Missing closing brace
	tmpfile, err := os.CreateTemp("", "invalid*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	loader := &Loader{}
	_, err = loader.Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoader_Load_EmptyFile(t *testing.T) {
	// Create an empty temporary file
	tmpfile, err := os.CreateTemp("", "empty*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	loader := &Loader{}
	_, err = loader.Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for empty JSON, got nil")
	}
}
func TestLoader_Save_Success(t *testing.T) {
	// Define sample AppConfig data
	configs := []AppConfig{
		{
			Name:        "TestApp",
			AllowedFrom: "08:00",
			AllowedTo:   "18:00",
		},
		{
			Name:        "TestApp2",
			AllowedFrom: "09:00",
			AllowedTo:   "20:00",
		},
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "output*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	loader := &Loader{}
	err = loader.Save(tmpfile.Name(), configs)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read the file back and verify its content
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	var savedConfigs []AppConfig
	err = json.Unmarshal(data, &savedConfigs)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved data: %v", err)
	}

	if len(savedConfigs) != len(configs) {
		t.Errorf("Expected %d configs, got %d", len(configs), len(savedConfigs))
	}

	for i, config := range savedConfigs {
		expected := configs[i]
		if config.Name != expected.Name {
			t.Errorf("Config %d: expected Name '%s', got '%s'", i, expected.Name, config.Name)
		}
		if config.AllowedFrom != expected.AllowedFrom {
			t.Errorf("Config %d: expected AllowedFrom '%s', got '%s'", i, expected.AllowedFrom, config.AllowedFrom)
		}
		if config.AllowedTo != expected.AllowedTo {
			t.Errorf("Config %d: expected AllowedTo '%s', got '%s'", i, expected.AllowedTo, config.AllowedTo)
		}
	}
}

func TestLoader_Save_InvalidPath(t *testing.T) {
	// Define sample AppConfig data
	configs := []AppConfig{
		{
			Name:        "TestApp",
			AllowedFrom: "08:00",
			AllowedTo:   "18:00",
		},
	}

	loader := &Loader{}
	// Attempt to save to an invalid path
	err := loader.Save("/invalid_path/config.json", configs)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
