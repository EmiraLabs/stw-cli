package infrastructure

import (
	"html/template"
	"io"
	"io/fs"
)

// FileSystem defines the interface for file operations
type FileSystem interface {
	WalkDir(root string, fn fs.WalkDirFunc) error
	ReadFile(filename string) ([]byte, error)
	Create(filename string) (io.WriteCloser, error)
	MkdirAll(path string, perm fs.FileMode) error
	RemoveAll(path string) error
}

// TemplateRenderer defines the interface for template rendering
type TemplateRenderer interface {
	ParseFiles(filenames ...string) (*template.Template, error)
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}
