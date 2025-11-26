package infrastructure

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// OSFileSystem implements FileSystem using os package
type OSFileSystem struct{}

// WalkDir walks the file tree rooted at root
func (fs *OSFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

// ReadFile reads the file named by filename and returns the contents
func (fs *OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// Create creates the named file with mode 0666 (before umask), truncating it if it already exists
func (fs *OSFileSystem) Create(filename string) (io.WriteCloser, error) {
	return os.Create(filename)
}

// MkdirAll creates a directory named path
func (fs *OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

// RemoveAll removes path and any children it contains
func (fs *OSFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// WriteFile writes data to a file named by filename
func (fs *OSFileSystem) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// Stat returns a FileInfo describing the named file
func (fs *OSFileSystem) Stat(filename string) (fs.FileInfo, error) {
	return os.Stat(filename)
}
