package meta

import (
	"testing"
)

func TestParseFrontMatter_JSON(t *testing.T) {
	jsonContent := `{"title": "JSON Title", "description": "JSON Description"}
<section>
<h1>Test</h1>
</section>`

	meta, body, err := ParseFrontMatter(jsonContent)
	if err != nil {
		t.Fatalf("ParseFrontMatter failed: %v", err)
	}
	if meta.Title != "JSON Title" {
		t.Errorf("Expected title 'JSON Title', got '%s'", meta.Title)
	}
	if meta.Description != "JSON Description" {
		t.Errorf("Expected description 'JSON Description', got '%s'", meta.Description)
	}
	expectedBody := `<section>
<h1>Test</h1>
</section>`
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestParseYAMLFrontMatter(t *testing.T) {
	// Valid YAML front matter
	yamlContent := `---
title: "YAML Title"
description: "YAML Description"
---
<body content>`

	meta, body, err := parseYAMLFrontMatter(yamlContent)
	if err != nil {
		t.Fatalf("parseYAMLFrontMatter failed: %v", err)
	}
	if meta.Title != "YAML Title" {
		t.Errorf("Expected title 'YAML Title', got '%s'", meta.Title)
	}
	if body != "<body content>" {
		t.Errorf("Expected body '<body content>', got '%s'", body)
	}

	// No YAML front matter
	noYAMLContent := `<html><body>Content</body></html>`
	meta2, body2, err2 := parseYAMLFrontMatter(noYAMLContent)
	if err2 != nil {
		t.Fatalf("parseYAMLFrontMatter failed on no YAML: %v", err2)
	}
	if meta2.Title != "" {
		t.Errorf("Expected empty title for no YAML, got '%s'", meta2.Title)
	}
	if body2 != noYAMLContent {
		t.Errorf("Expected unchanged body, got '%s'", body2)
	}

	// Invalid YAML front matter (missing closing ---)
	invalidYAML := `---
title: "Test"
<body>`
	_, _, err3 := parseYAMLFrontMatter(invalidYAML)
	if err3 == nil {
		t.Error("Expected error for invalid YAML front matter")
	}
}

func TestParseJSONFrontMatter(t *testing.T) {
	// Valid JSON front matter
	jsonContent := `{"title": "JSON Title", "description": "JSON Description"}
<body content>`

	meta, body, err := parseJSONFrontMatter(jsonContent)
	if err != nil {
		t.Fatalf("parseJSONFrontMatter failed: %v", err)
	}
	if meta.Title != "JSON Title" {
		t.Errorf("Expected title 'JSON Title', got '%s'", meta.Title)
	}
	if body != "<body content>" {
		t.Errorf("Expected body '<body content>', got '%s'", body)
	}

	// No JSON front matter
	noJSONContent := `<html><body>Content</body></html>`
	meta2, body2, err2 := parseJSONFrontMatter(noJSONContent)
	if err2 != nil {
		t.Fatalf("parseJSONFrontMatter failed on no JSON: %v", err2)
	}
	if meta2.Title != "" {
		t.Errorf("Expected empty title for no JSON, got '%s'", meta2.Title)
	}
	if body2 != noJSONContent {
		t.Errorf("Expected unchanged body, got '%s'", body2)
	}

	// Invalid JSON front matter (missing closing })
	invalidJSON := `{"title": "Test"
<body>`
	_, _, err3 := parseJSONFrontMatter(invalidJSON)
	if err3 == nil {
		t.Error("Expected error for invalid JSON front matter")
	}
}

func TestMerge(t *testing.T) {
	siteMeta := Meta{
		Title:       "Site Title",
		Description: "Site Description",
		Robots:      "index,follow",
	}
	pageMeta := Meta{
		Title:       "Page Title",
		Description: "",
		Keywords:    "page,keywords",
	}

	merged := Merge(siteMeta, pageMeta)
	if merged.Title != "Page Title" {
		t.Errorf("Expected title 'Page Title', got '%s'", merged.Title)
	}
	if merged.Description != "Site Description" {
		t.Errorf("Expected description 'Site Description', got '%s'", merged.Description)
	}
	if merged.Keywords != "page,keywords" {
		t.Errorf("Expected keywords 'page,keywords', got '%s'", merged.Keywords)
	}
	if merged.Robots != "index,follow" {
		t.Errorf("Expected robots 'index,follow', got '%s'", merged.Robots)
	}
}

func TestLoadSiteMeta(t *testing.T) {
	config := map[string]interface{}{
		"meta": map[string]interface{}{
			"title":       "Site Title",
			"description": "Site Description",
			"keywords":    "site, keywords",
		},
		"other": "value",
	}

	meta := LoadSiteMeta(config)
	if meta.Title != "Site Title" {
		t.Errorf("Expected title 'Site Title', got '%s'", meta.Title)
	}
	if meta.Description != "Site Description" {
		t.Errorf("Expected description 'Site Description', got '%s'", meta.Description)
	}
	if meta.Keywords != "site, keywords" {
		t.Errorf("Expected keywords 'site, keywords', got '%s'", meta.Keywords)
	}
}

func TestLoadSiteMeta_NoMeta(t *testing.T) {
	config := map[string]interface{}{
		"other": "value",
	}

	meta := LoadSiteMeta(config)
	if meta.Title != "" {
		t.Errorf("Expected empty title, got '%s'", meta.Title)
	}
}

func TestValidate(t *testing.T) {
	meta := Meta{
		Title:       "This is a very long title that exceeds sixty characters in length",
		Description: "Short desc",
		OgImage:     "/assets/image.jpg",
	}
	err := meta.Validate("assets")
	if err == nil {
		t.Error("Expected validation error for long title")
	}

	meta.Title = "Short Title"
	err = meta.Validate("assets")
	if err != nil {
		t.Errorf("Unexpected validation error: %v", err)
	}

	meta.OgImage = "/other/image.jpg"
	err = meta.Validate("assets")
	if err == nil {
		t.Error("Expected validation error for og_image not under /assets/")
	}
}
