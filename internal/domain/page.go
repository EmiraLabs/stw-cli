package domain

import (
	"html/template"

	"github.com/EmiraLabs/stw-cli/internal/meta"
)

// Page represents a web page with title and content
type Page struct {
	Title   string
	Content template.HTML
	Path    string // relative path
	IsDev   bool
	Config  map[string]interface{}
	Meta    meta.Meta
}
