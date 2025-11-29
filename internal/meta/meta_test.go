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

	// JSON front matter ending with } (no newline)
	jsonContentNoNewline := `{"title": "JSON Title No NL", "description": "JSON Description No NL"}<body content without newline>`
	meta3, body3, err3 := parseJSONFrontMatter(jsonContentNoNewline)
	if err3 != nil {
		t.Fatalf("parseJSONFrontMatter failed on JSON without newline: %v", err3)
	}
	if meta3.Title != "JSON Title No NL" {
		t.Errorf("Expected title 'JSON Title No NL', got '%s'", meta3.Title)
	}
	if body3 != "<body content without newline>" {
		t.Errorf("Expected body '<body content without newline>', got '%s'", body3)
	}

	// Invalid JSON front matter (missing closing })
	invalidJSON := `{"title": "Test"
<body>`
	_, _, err4 := parseJSONFrontMatter(invalidJSON)
	if err4 == nil {
		t.Error("Expected error for invalid JSON front matter")
	}
}

func TestParseFrontMatter(t *testing.T) {
	// Test YAML front matter
	yamlContent := `---
title: "YAML Title"
description: "YAML Description"
---
<body>Content</body>`

	meta, body, err := ParseFrontMatter(yamlContent)
	if err != nil {
		t.Fatalf("ParseFrontMatter failed on YAML: %v", err)
	}
	if meta.Title != "YAML Title" {
		t.Errorf("Expected title 'YAML Title', got '%s'", meta.Title)
	}
	if meta.Description != "YAML Description" {
		t.Errorf("Expected description 'YAML Description', got '%s'", meta.Description)
	}
	if body != "<body>Content</body>" {
		t.Errorf("Expected body '<body>Content</body>', got '%s'", body)
	}

	// Test JSON front matter
	jsonContent := `{"title": "JSON Title", "description": "JSON Description"}
<body>Content</body>`

	meta2, body2, err2 := ParseFrontMatter(jsonContent)
	if err2 != nil {
		t.Fatalf("ParseFrontMatter failed on JSON: %v", err2)
	}
	if meta2.Title != "JSON Title" {
		t.Errorf("Expected title 'JSON Title', got '%s'", meta2.Title)
	}
	if meta2.Description != "JSON Description" {
		t.Errorf("Expected description 'JSON Description', got '%s'", meta2.Description)
	}
	if body2 != "<body>Content</body>" {
		t.Errorf("Expected body '<body>Content</body>', got '%s'", body2)
	}
	// Test JSON front matter starting with {\n
	jsonContent2 := `{
  "title": "JSON Title 2",
  "description": "JSON Description 2"
}
<body>Content 2</body>`

	meta3, body3, err3 := parseJSONFrontMatter(jsonContent2)
	if err3 != nil {
		t.Fatalf("parseJSONFrontMatter failed on JSON with newline: %v", err3)
	}
	if meta3.Title != "JSON Title 2" {
		t.Errorf("Expected title 'JSON Title 2', got '%s'", meta3.Title)
	}
	if body3 != "<body>Content 2</body>" {
		t.Errorf("Expected body '<body>Content 2</body>', got '%s'", body3)
	}

	// Test empty YAML front matter (should fall through to JSON)
	emptyYAMLContent := `---
---
<body>Content</body>`

	meta4, body4, err4 := ParseFrontMatter(emptyYAMLContent)
	if err4 != nil {
		t.Fatalf("ParseFrontMatter failed on empty YAML: %v", err4)
	}
	// Should have no meta
	if meta4.Title != "" {
		t.Errorf("Expected empty title, got '%s'", meta4.Title)
	}
	if body4 != emptyYAMLContent {
		t.Errorf("Expected unchanged body, got '%s'", body4)
	}

	// Test no front matter
	noFMContent := `<html><body>Content</body></html>`
	meta5, body5, err5 := ParseFrontMatter(noFMContent)
	if err5 != nil {
		t.Fatalf("ParseFrontMatter failed on no front matter: %v", err5)
	}
	if meta5.Title != "" {
		t.Errorf("Expected empty title, got '%s'", meta5.Title)
	}
	if body5 != noFMContent {
		t.Errorf("Expected unchanged body, got '%s'", body5)
	}

	// Test invalid YAML (should return error)
	invalidYAML := `---
title: "Test"
invalid: yaml: content: [
---
<body>`

	_, _, err6 := ParseFrontMatter(invalidYAML)
	if err6 == nil {
		t.Error("Expected error for invalid YAML front matter")
	}

	// Test invalid JSON (should return error)
	invalidJSON := `{"title": "Test"
<body>`

	_, _, err7 := ParseFrontMatter(invalidJSON)
	if err7 == nil {
		t.Error("Expected error for invalid JSON front matter")
	}
}

func TestMerge(t *testing.T) {
	siteMeta := Meta{
		Title:              "Site Title",
		Description:        "Site Description",
		Canonical:          "https://example.com",
		Robots:             "index,follow",
		Keywords:           "site,keywords",
		OgTitle:            "Site OG Title",
		OgDescription:      "Site OG Description",
		OgImage:            "/assets/site-og.jpg",
		TwitterTitle:       "Site Twitter Title",
		TwitterDescription: "Site Twitter Description",
		TwitterImage:       "/assets/site-twitter.jpg",
		JsonLd:             map[string]interface{}{"@type": "Organization"},
	}
	pageMeta := Meta{
		Title:              "Page Title",
		Description:        "Page Description",
		Canonical:          "https://example.com/page",
		Robots:             "noindex",
		Keywords:           "page,keywords",
		OgTitle:            "Page OG Title",
		OgDescription:      "Page OG Description",
		OgImage:            "/assets/page-og.jpg",
		TwitterTitle:       "Page Twitter Title",
		TwitterDescription: "Page Twitter Description",
		TwitterImage:       "/assets/page-twitter.jpg",
		JsonLd:             map[string]interface{}{"@type": "Article"},
	}

	merged := Merge(siteMeta, pageMeta)
	// All pageMeta fields should override siteMeta
	if merged.Title != "Page Title" {
		t.Errorf("Expected title 'Page Title', got '%s'", merged.Title)
	}
	if merged.Description != "Page Description" {
		t.Errorf("Expected description 'Page Description', got '%s'", merged.Description)
	}
	if merged.Canonical != "https://example.com/page" {
		t.Errorf("Expected canonical 'https://example.com/page', got '%s'", merged.Canonical)
	}
	if merged.Robots != "noindex" {
		t.Errorf("Expected robots 'noindex', got '%s'", merged.Robots)
	}
	if merged.Keywords != "page,keywords" {
		t.Errorf("Expected keywords 'page,keywords', got '%s'", merged.Keywords)
	}
	if merged.OgTitle != "Page OG Title" {
		t.Errorf("Expected og_title 'Page OG Title', got '%s'", merged.OgTitle)
	}
	if merged.OgDescription != "Page OG Description" {
		t.Errorf("Expected og_description 'Page OG Description', got '%s'", merged.OgDescription)
	}
	if merged.OgImage != "/assets/page-og.jpg" {
		t.Errorf("Expected og_image '/assets/page-og.jpg', got '%s'", merged.OgImage)
	}
	if merged.TwitterTitle != "Page Twitter Title" {
		t.Errorf("Expected twitter_title 'Page Twitter Title', got '%s'", merged.TwitterTitle)
	}
	if merged.TwitterDescription != "Page Twitter Description" {
		t.Errorf("Expected twitter_description 'Page Twitter Description', got '%s'", merged.TwitterDescription)
	}
	if merged.TwitterImage != "/assets/page-twitter.jpg" {
		t.Errorf("Expected twitter_image '/assets/page-twitter.jpg', got '%s'", merged.TwitterImage)
	}
	if merged.JsonLd == nil || merged.JsonLd["@type"] != "Article" {
		t.Errorf("Expected JsonLd to be overridden, got %v", merged.JsonLd)
	}
}

func TestMerge_EmptyPageMeta(t *testing.T) {
	siteMeta := Meta{
		Title:       "Site Title",
		Description: "Site Description",
		Keywords:    "site,keywords",
	}
	pageMeta := Meta{} // Empty

	merged := Merge(siteMeta, pageMeta)
	// Should keep all siteMeta values
	if merged.Title != "Site Title" {
		t.Errorf("Expected title 'Site Title', got '%s'", merged.Title)
	}
	if merged.Description != "Site Description" {
		t.Errorf("Expected description 'Site Description', got '%s'", merged.Description)
	}
	if merged.Keywords != "site,keywords" {
		t.Errorf("Expected keywords 'site,keywords', got '%s'", merged.Keywords)
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
	meta.Description = "This is a very long description that exceeds one hundred and sixty characters in total length and should trigger a validation error when the Validate method is called with this input string value."
	err = meta.Validate("assets")
	if err == nil {
		t.Error("Expected validation error for long description")
	}

	meta.Description = "Short desc"
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
