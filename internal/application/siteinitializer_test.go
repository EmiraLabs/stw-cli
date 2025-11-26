package application

import (
	"testing"
)

func TestSiteInitializer_InitTailwind(t *testing.T) {
	fs := NewMockFileSystem()
	initializer := NewSiteInitializer(fs)

	err := initializer.InitTailwind()
	if err != nil {
		t.Fatalf("InitTailwind failed: %v", err)
	}

	// Check if package.json was created
	if _, err := fs.Stat("package.json"); err != nil {
		t.Errorf("package.json not created")
	}

	// Check if postcss.config.js was created
	if _, err := fs.Stat("postcss.config.js"); err != nil {
		t.Errorf("postcss.config.js not created")
	}

	// Check if tailwind.config.js was created
	if _, err := fs.Stat("tailwind.config.js"); err != nil {
		t.Errorf("tailwind.config.js not created")
	}

	// Check if styles.css was updated
	if _, err := fs.Stat("assets/css/styles.css"); err != nil {
		t.Errorf("styles.css not updated")
	}

	// Check content of styles.css
	content, err := fs.ReadFile("assets/css/styles.css")
	if err != nil {
		t.Errorf("Failed to read styles.css: %v", err)
	}
	expected := "@tailwind base;\n@tailwind components;\n@tailwind utilities;"
	if string(content) != expected {
		t.Errorf("styles.css content mismatch. Got: %s, Expected: %s", string(content), expected)
	}
}
