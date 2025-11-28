package main

import (
	"html/template"
	"os"
	"testing"
)

func TestConvertToHTML(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  interface{}
		checkFunc func(interface{}) bool
	}{
		{
			name:     "string conversion",
			input:    "test",
			expected: template.HTML("test"),
		},
		{
			name:     "HTML string conversion",
			input:    "<h1>test</h1>",
			expected: template.HTML("<h1>test</h1>"),
		},
		{
			name: "map conversion",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"nested": "nested_value",
				},
			},
			checkFunc: func(result interface{}) bool {
				if m, ok := result.(map[string]interface{}); ok {
					if v1, ok := m["key1"].(template.HTML); ok && string(v1) == "value1" {
						if m2, ok := m["key2"].(map[string]interface{}); ok {
							if v2, ok := m2["nested"].(template.HTML); ok && string(v2) == "nested_value" {
								return true
							}
						}
					}
				}
				return false
			},
		},
		{
			name: "array conversion",
			input: []interface{}{
				"item1",
				map[string]interface{}{
					"key": "value",
				},
			},
			checkFunc: func(result interface{}) bool {
				if arr, ok := result.([]interface{}); ok && len(arr) == 2 {
					if v1, ok := arr[0].(template.HTML); ok && string(v1) == "item1" {
						if m, ok := arr[1].(map[string]interface{}); ok {
							if v2, ok := m["key"].(template.HTML); ok && string(v2) == "value" {
								return true
							}
						}
					}
				}
				return false
			},
		},
		{
			name:     "non-string value",
			input:    42,
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToHTML(tt.input)
			if tt.checkFunc != nil {
				if !tt.checkFunc(result) {
					t.Errorf("convertToHTML() result does not match expected structure")
				}
			} else if result != tt.expected {
				t.Errorf("convertToHTML() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory and config file
	tmpDir, err := os.MkdirTemp("", "stw-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configContent := `
navigations:
  - title: "Home"
    url: "/"
  - title: "About"
    url: "/about"
home:
  title: "<h1>Welcome</h1>"
  content: "Some content"
`
	configPath := tmpDir + "/config.yaml"
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test loading config
	config, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}

	// Check that navigations exist
	if config["navigations"] == nil {
		t.Error("navigations not found in config")
	}

	// Check that home.title is HTML
	if home, ok := config["home"].(map[string]interface{}); ok {
		if title, ok := home["title"].(template.HTML); ok {
			if string(title) != "<h1>Welcome</h1>" {
				t.Errorf("Expected HTML title, got %v", title)
			}
		} else {
			t.Error("title is not template.HTML")
		}
	} else {
		t.Error("home is not a map")
	}
}

func TestLoadConfigNoFile(t *testing.T) {
	// Test loading config when file doesn't exist
	// Change to a temp dir
	tmpDir, err := os.MkdirTemp("", "stw-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	config, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if len(config) != 0 {
		t.Errorf("Expected empty config, got %v", config)
	}
}
