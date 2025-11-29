package infrastructure

import (
	"bytes"
	"html/template"
	"os"
	"testing"
)

func TestGoTemplateRenderer_ParseFiles(t *testing.T) {
	renderer := &GoTemplateRenderer{}
	tempDir := t.TempDir()
	// Create template files
	baseContent := `<html>{{.Title}}</html>`
	headerContent := `<head></head>`
	footerContent := `<footer></footer>`

	baseFile := tempDir + "/base.html"
	headerFile := tempDir + "/header.html"
	footerFile := tempDir + "/footer.html"

	os.WriteFile(baseFile, []byte(baseContent), 0644)
	os.WriteFile(headerFile, []byte(headerContent), 0644)
	os.WriteFile(footerFile, []byte(footerContent), 0644)

	tmpl, err := renderer.ParseFiles(baseFile, headerFile, footerFile)
	if err != nil {
		t.Errorf("ParseFiles failed: %v", err)
	}
	if tmpl == nil {
		t.Error("Template is nil")
	}
}

func TestGoTemplateRenderer_ParseFiles_Error(t *testing.T) {
	renderer := &GoTemplateRenderer{}
	
	// Test file not found error
	_, err := renderer.ParseFiles("nonexistent.html")
	if err == nil {
		t.Error("Expected error from ParseFiles for non-existent file, got nil")
	}

	// Test template parse error (invalid syntax)
	tempDir := t.TempDir()
	invalidTemplateFile := tempDir + "/invalid.html"
	invalidContent := `{{.Title`  // Missing closing braces
	os.WriteFile(invalidTemplateFile, []byte(invalidContent), 0644)
	
	_, err = renderer.ParseFiles(invalidTemplateFile)
	if err == nil {
		t.Error("Expected error from ParseFiles for invalid template syntax, got nil")
	}
}

func TestGoTemplateRenderer_ExecuteTemplate(t *testing.T) {
	renderer := &GoTemplateRenderer{}
	tmpl, _ := template.New("test").Parse("{{.Title}}")
	renderer.tmpl = tmpl

	var buf bytes.Buffer
	data := struct{ Title string }{"Test Title"}
	err := renderer.ExecuteTemplate(&buf, "test", data)
	if err != nil {
		t.Errorf("ExecuteTemplate failed: %v", err)
	}
	if buf.String() != "Test Title" {
		t.Errorf("Expected 'Test Title', got '%s'", buf.String())
	}
}

func TestGoTemplateRenderer_ExecuteTemplate_Error(t *testing.T) {
	renderer := &GoTemplateRenderer{}
	tmpl, _ := template.New("test").Parse("{{.Invalid}}")
	renderer.tmpl = tmpl

	var buf bytes.Buffer
	data := struct{ Title string }{"Test Title"}
	err := renderer.ExecuteTemplate(&buf, "test", data)
	if err == nil {
		t.Errorf("Expected error from ExecuteTemplate, but got: %v", err)
	}
}
