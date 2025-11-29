package main

import (
	"html/template"
	"os"
	"path/filepath"
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

func TestLoadConfigInvalidYAML(t *testing.T) {
	// Create a temporary directory and invalid config file
	tmpDir, err := os.MkdirTemp("", "stw-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	invalidYAML := `
navigations:
  - title: "Home"
    url: "/"
  - title: "About"
    url: "/about"
    invalid: [unclosed
`
	configPath := tmpDir + "/config.yaml"
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test loading config
	_, err = loadConfig()
	if err == nil {
		t.Error("Expected error for invalid YAML, but got none")
	}
}

func TestLoadConfigValid(t *testing.T) {
	// Create a temporary directory and valid config file
	tmpDir, err := os.MkdirTemp("", "stw-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	validYAML := `
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
	if err := os.WriteFile(configPath, []byte(validYAML), 0644); err != nil {
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

	if config["navigations"] == nil {
		t.Error("navigations not found in config")
	}

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

func TestBuild(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "stw-build-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create pages/index.html
	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	indexPath := filepath.Join(pagesDir, "index.html")
	content := `---
title: Home
---
<h1>Welcome</h1>`
	if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create templates
	templatesDir := filepath.Join(tmpDir, "templates")
	os.MkdirAll(filepath.Join(templatesDir, "components"), 0755)
	os.MkdirAll(filepath.Join(templatesDir, "partials"), 0755)

	basePath := filepath.Join(templatesDir, "base.html")
	baseContent := `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
<header><h1>{{.Title}}</h1></header>
{{.Content}}
<footer>Footer</footer>
</body>
</html>`
	if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	headerPath := filepath.Join(templatesDir, "components", "header.html")
	if err := os.WriteFile(headerPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	footerPath := filepath.Join(templatesDir, "components", "footer.html")
	if err := os.WriteFile(footerPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	headPath := filepath.Join(templatesDir, "partials", "head.html")
	if err := os.WriteFile(headPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Create assets
	assetsDir := filepath.Join(tmpDir, "assets")
	os.MkdirAll(assetsDir, 0755)
	cssDir := filepath.Join(assetsDir, "css")
	os.MkdirAll(cssDir, 0755)
	cssPath := filepath.Join(cssDir, "styles.css")
	if err := os.WriteFile(cssPath, []byte("body { }"), 0644); err != nil {
		t.Fatal(err)
	}

	// Change wd
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	// Call build
	if err := build(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	// Check dist/index.html
	distPath := filepath.Join(tmpDir, "dist", "index.html")
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		t.Error("dist/index.html not created")
	}

	// Optionally check content
	data, err := os.ReadFile(distPath)
	if err != nil {
		t.Fatal(err)
	}
	expected := `<!DOCTYPE html>
<html>
<head><title>Home</title></head>
<body>
<header><h1>Home</h1></header>
<h1>Welcome</h1>
<footer>Footer</footer>
</body>
</html>`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	// Create a temp directory without config.yaml
	tmpDir, err := os.MkdirTemp("", "stw-test-missing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test loading config when file doesn't exist
	config, err := loadConfig()
	if err != nil {
		t.Fatalf("Expected no error for missing config, got: %v", err)
	}

	// Should return empty config map
	if config == nil {
		t.Error("Expected empty config map, got nil")
	}
	if len(config) != 0 {
		t.Errorf("Expected empty config, got %d items", len(config))
	}
}

func TestBuild_ConfigError(t *testing.T) {
	// Create temp dir with invalid YAML config
	tmpDir, err := os.MkdirTemp("", "stw-build-error-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create invalid config
	invalidYAML := `
navigations:
  - title: "Home"
    url: "/"
    invalid: [unclosed
`
	configPath := tmpDir + "/config.yaml"
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	// Call build - should fail due to invalid YAML
	err = build()
	if err == nil {
		t.Error("Expected error from build with invalid config, got nil")
	}
}

