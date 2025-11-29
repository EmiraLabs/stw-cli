package domain

import (
	"testing"

	"github.com/EmiraLabs/stw-cli/internal/meta"
)

func TestPageStruct(t *testing.T) {
	p := Page{
		Title:   "Test Page",
		Content: "<p>Content</p>",
		Path:    "/test",
		IsDev:   true,
		Config:  map[string]interface{}{"key": "value"},
		Meta:    meta.Meta{},
	}

	if p.Title != "Test Page" {
		t.Errorf("Expected Title 'Test Page', got %s", p.Title)
	}

	if p.Path != "/test" {
		t.Errorf("Expected Path '/test', got %s", p.Path)
	}

	if !p.IsDev {
		t.Error("Expected IsDev true")
	}
}
