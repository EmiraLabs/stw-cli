package meta

import (
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	yamlContent := `---
title: "Test Title"
description: "Test Description"
---
<section>
<h1>Test</h1>
</section>`

	meta, body, err := ParseFrontMatter(yamlContent)
	if err != nil {
		t.Fatalf("ParseFrontMatter failed: %v", err)
	}
	if meta.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", meta.Title)
	}
	if meta.Description != "Test Description" {
		t.Errorf("Expected description 'Test Description', got '%s'", meta.Description)
	}
	expectedBody := `<section>
<h1>Test</h1>
</section>`
	if body != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
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
