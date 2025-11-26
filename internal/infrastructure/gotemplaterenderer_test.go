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
