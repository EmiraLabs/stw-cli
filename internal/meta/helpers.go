// Package meta provides functionality for parsing, validating, and merging SEO metadata
// for static site pages. It supports YAML front matter in page files and site-wide
// defaults from configuration.
package meta

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Meta represents SEO metadata for a page, including standard meta tags,
// Open Graph properties, Twitter Card data, and JSON-LD structured data.
type Meta struct {
	Title              string                 `yaml:"title" json:"title"`
	Description        string                 `yaml:"description" json:"description"`
	Canonical          string                 `yaml:"canonical" json:"canonical"`
	Robots             string                 `yaml:"robots" json:"robots"`
	Keywords           string                 `yaml:"keywords" json:"keywords"`
	OgTitle            string                 `yaml:"og_title" json:"og_title"`
	OgDescription      string                 `yaml:"og_description" json:"og_description"`
	OgImage            string                 `yaml:"og_image" json:"og_image"`
	TwitterTitle       string                 `yaml:"twitter_title" json:"twitter_title"`
	TwitterDescription string                 `yaml:"twitter_description" json:"twitter_description"`
	TwitterImage       string                 `yaml:"twitter_image" json:"twitter_image"`
	JsonLd             map[string]interface{} `yaml:"jsonld" json:"jsonld"`
}

// Validate checks the meta fields for SEO best practices and constraints.
// It ensures title and description lengths are within recommended limits,
// and that Open Graph images are properly located under the assets directory.
func (m *Meta) Validate(assetsDir string) error {
	if len(m.Title) > 60 {
		return fmt.Errorf("title exceeds 60 characters: %s", m.Title)
	}
	if len(m.Description) > 160 {
		return fmt.Errorf("description exceeds 160 characters: %s", m.Description)
	}
	if m.OgImage != "" {
		// Check if og_image exists under assets/
		if !strings.HasPrefix(m.OgImage, "/assets/") {
			return fmt.Errorf("og_image must be under /assets/: %s", m.OgImage)
		}
		// For now, assume it exists; in build, we can check file existence
	}
	return nil
}

// Merge combines site-level default meta with page-specific overrides.
// Page meta values take precedence when present (non-empty).
func Merge(siteMeta, pageMeta Meta) Meta {
	merged := siteMeta
	if pageMeta.Title != "" {
		merged.Title = pageMeta.Title
	}
	if pageMeta.Description != "" {
		merged.Description = pageMeta.Description
	}
	if pageMeta.Canonical != "" {
		merged.Canonical = pageMeta.Canonical
	}
	if pageMeta.Robots != "" {
		merged.Robots = pageMeta.Robots
	}
	if pageMeta.Keywords != "" {
		merged.Keywords = pageMeta.Keywords
	}
	if pageMeta.OgTitle != "" {
		merged.OgTitle = pageMeta.OgTitle
	}
	if pageMeta.OgDescription != "" {
		merged.OgDescription = pageMeta.OgDescription
	}
	if pageMeta.OgImage != "" {
		merged.OgImage = pageMeta.OgImage
	}
	if pageMeta.TwitterTitle != "" {
		merged.TwitterTitle = pageMeta.TwitterTitle
	}
	if pageMeta.TwitterDescription != "" {
		merged.TwitterDescription = pageMeta.TwitterDescription
	}
	if pageMeta.TwitterImage != "" {
		merged.TwitterImage = pageMeta.TwitterImage
	}
	if len(pageMeta.JsonLd) > 0 {
		merged.JsonLd = pageMeta.JsonLd
	}
	return merged
}

// ParseFrontMatter extracts YAML front matter from page content.
// It looks for content between --- markers and parses it into a Meta struct.
// Returns the parsed meta, the body content without front matter, and any error.
func ParseFrontMatter(content string) (Meta, string, error) {
	var meta Meta
	var body string

	// Check for YAML front matter
	if strings.HasPrefix(content, "---\n") {
		parts := strings.SplitN(content, "---\n", 3)
		if len(parts) >= 3 {
			if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
				return meta, content, err
			}
			body = parts[2]
		} else {
			body = content
		}
	} else if strings.HasPrefix(content, "{\n") || strings.HasPrefix(content, "{") {
		// Simple check for JSON
		end := strings.Index(content, "}\n")
		if end == -1 {
			end = strings.Index(content, "}")
		}
		if end != -1 {
			jsonPart := content[:end+1]
			if err := json.Unmarshal([]byte(jsonPart), &meta); err != nil {
				return meta, content, err
			}
			body = content[end+1:]
		} else {
			body = content
		}
	} else {
		body = content
	}

	return meta, body, nil
}

// LoadSiteMeta extracts site-wide meta configuration from the config map.
// It looks for a "meta" key in the configuration and unmarshals it into a Meta struct.
func LoadSiteMeta(config map[string]interface{}) Meta {
	var meta Meta
	if metaData, ok := config["meta"].(map[string]interface{}); ok {
		// Convert map to Meta
		data, _ := yaml.Marshal(metaData)
		yaml.Unmarshal(data, &meta)
	}
	return meta
}
