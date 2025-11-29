package application

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/EmiraLabs/stw-cli/internal/meta"
)

// MockFileSystem is a mock implementation of FileSystem for testing
type MockFileSystem struct {
	files       map[string][]byte
	dirs        map[string]bool
	walkCalls   []string
	createCalls []string
	mkdirCalls  []string
	removeCalls []string
	readError   error
	createError error
	mkdirError  error
	removeError error
	walkError   error
	hasPages    bool
	hasAssets   bool
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files:       make(map[string][]byte),
		dirs:        make(map[string]bool),
		walkCalls:   []string{},
		createCalls: []string{},
		mkdirCalls:  []string{},
		removeCalls: []string{},
		hasPages:    true,
		hasAssets:   true,
	}
}

func (m *MockFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	if m.walkError != nil {
		return m.walkError
	}
	m.walkCalls = append(m.walkCalls, root)
	// Simulate walking pages
	if root == "pages" && m.hasPages {
		// index.html
		if err := fn("pages/index.html", &mockDirEntry{name: domain.IndexFile, isDir: false}, nil); err != nil {
			return err
		}
		// about/index.html
		if err := fn("pages/about/index.html", &mockDirEntry{name: domain.IndexFile, isDir: false}, nil); err != nil {
			return err
		}
		// about/contact/index.html
		if err := fn("pages/about/contact/index.html", &mockDirEntry{name: domain.IndexFile, isDir: false}, nil); err != nil {
			return err
		}
		// some other file
		if err := fn("pages/other.txt", &mockDirEntry{name: "other.txt", isDir: false}, nil); err != nil {
			return err
		}
	}
	if root == "assets" && m.hasAssets {
		if err := fn("assets", &mockDirEntry{name: "assets", isDir: true}, nil); err != nil {
			return err
		}
		if err := fn("assets/css", &mockDirEntry{name: "css", isDir: true}, nil); err != nil {
			return err
		}
		if err := fn("assets/css/style.css", &mockDirEntry{name: "style.css", isDir: false}, nil); err != nil {
			return err
		}
		if err := fn("assets/js", &mockDirEntry{name: "js", isDir: true}, nil); err != nil {
			return err
		}
		if err := fn("assets/js/app.js", &mockDirEntry{name: "app.js", isDir: false}, nil); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	if m.readError != nil {
		return nil, m.readError
	}
	if content, ok := m.files[filename]; ok {
		return content, nil
	}
	// Default content for pages
	if strings.HasSuffix(filename, domain.IndexFile) {
		return []byte("<h1>{{.Title}}</h1>{{.Config.test}}"), nil
	}
	return []byte("content"), nil
}

func (m *MockFileSystem) Create(filename string) (io.WriteCloser, error) {
	m.createCalls = append(m.createCalls, filename)
	if m.createError != nil {
		return nil, m.createError
	}
	return &mockWriteCloser{buffer: &bytes.Buffer{}}, nil
}

func (m *MockFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	m.mkdirCalls = append(m.mkdirCalls, path)
	if m.mkdirError != nil {
		return m.mkdirError
	}
	m.dirs[path] = true
	return nil
}

func (m *MockFileSystem) RemoveAll(path string) error {
	m.removeCalls = append(m.removeCalls, path)
	if m.removeError != nil {
		return m.removeError
	}
	return nil
}

type mockDirEntry struct {
	name  string
	isDir bool
}

func (m *mockDirEntry) Name() string               { return m.name }
func (m *mockDirEntry) IsDir() bool                { return m.isDir }
func (m *mockDirEntry) Type() fs.FileMode          { return 0 }
func (m *mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

type mockWriteCloser struct {
	buffer *bytes.Buffer
}

func (m *mockWriteCloser) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *mockWriteCloser) Close() error {
	return nil
}

// MockTemplateRenderer is a mock implementation of TemplateRenderer
type MockTemplateRenderer struct {
	parseFilesCalls [][]string
	executeCalls    []executeCall
	tmpl            *template.Template
	parseError      error
	executeError    error
}

type executeCall struct {
	name string
	data interface{}
}

func NewMockTemplateRenderer() *MockTemplateRenderer {
	return &MockTemplateRenderer{
		parseFilesCalls: [][]string{},
		executeCalls:    []executeCall{},
	}
}

func (m *MockTemplateRenderer) ParseFiles(filenames ...string) (*template.Template, error) {
	m.parseFilesCalls = append(m.parseFilesCalls, filenames)
	if m.parseError != nil {
		return nil, m.parseError
	}
	// Return a simple template
	var tmplStr string
	if m.executeError != nil {
		tmplStr = `{{.Invalid}}`
	} else {
		tmplStr = `{{.Title}}{{.Content}}`
	}
	tmpl, _ := template.New("base.html").Parse(tmplStr)
	m.tmpl = tmpl
	return tmpl, nil
}

func (m *MockTemplateRenderer) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	m.executeCalls = append(m.executeCalls, executeCall{name: name, data: data})
	if m.executeError != nil {
		return m.executeError
	}
	return m.tmpl.ExecuteTemplate(wr, name, data)
}

func TestNewSiteBuilder(t *testing.T) {
	site := &domain.Site{}
	fs := NewMockFileSystem()
	renderer := NewMockTemplateRenderer()
	builder := NewSiteBuilder(site, fs, renderer)
	if builder.site != site || builder.fs != fs || builder.renderer != renderer {
		t.Error("NewSiteBuilder did not set fields correctly")
	}
}

func TestSiteBuilder_Build(t *testing.T) {
	tests := []struct {
		name        string
		setupFS     func(*MockFileSystem)
		setupRender func(*MockTemplateRenderer)
		expectError bool
		errorMsg    string
		checkFunc   func(*testing.T, *MockFileSystem, *MockTemplateRenderer)
	}{
		{
			name:        "Success",
			setupFS:     func(fs *MockFileSystem) {},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: false,
			checkFunc: func(t *testing.T, fs *MockFileSystem, r *MockTemplateRenderer) {
				// Check that RemoveAll was called
				if len(fs.removeCalls) != 1 || fs.removeCalls[0] != "dist" {
					t.Error("RemoveAll not called correctly")
				}
				// Check MkdirAll for dist
				if len(fs.mkdirCalls) < 1 || fs.mkdirCalls[0] != "dist" {
					t.Error("MkdirAll for dist not called")
				}
				// Check ParseFiles
				if len(r.parseFilesCalls) != 1 {
					t.Error("ParseFiles not called")
				}
				// Check that pages were built
				expectedCreates := []string{"dist/index.html", "dist/about/index.html", "dist/about/contact/index.html"}
				for _, expected := range expectedCreates {
					if !contains(fs.createCalls, expected) {
						t.Errorf("Expected create call for %s", expected)
					}
				}
				// Check assets copied
				expectedAssetCreates := []string{"dist/assets/css/style.css", "dist/assets/js/app.js"}
				for _, expected := range expectedAssetCreates {
					if !contains(fs.createCalls, expected) {
						t.Errorf("Expected create call for %s", expected)
					}
				}
			},
		},
		{
			name: "Parse Error",
			setupFS: func(fs *MockFileSystem) {},
			setupRender: func(r *MockTemplateRenderer) {
				r.parseError = errors.New("parse error")
			},
			expectError: true,
			errorMsg:    "parse error",
		},
		{
			name: "Mkdir Error",
			setupFS: func(fs *MockFileSystem) {
				fs.mkdirError = errors.New("mkdir error")
			},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: true,
			errorMsg:    "mkdir error",
		},
		{
			name: "Remove Error",
			setupFS: func(fs *MockFileSystem) {
				fs.removeError = errors.New("remove error")
			},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: true,
			errorMsg:    "remove error",
		},
		{
			name: "No Pages",
			setupFS: func(fs *MockFileSystem) {
				fs.hasPages = false
			},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: false,
			checkFunc: func(t *testing.T, fs *MockFileSystem, r *MockTemplateRenderer) {
				// Check that no pages were built
				for _, call := range fs.createCalls {
					if strings.HasPrefix(call, "dist/") && strings.HasSuffix(call, ".html") {
						t.Errorf("Unexpected create call for page: %s", call)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := &domain.Site{
				PagesDir:         "pages",
				TemplatesDir:     "templates",
				AssetsDir:        "assets",
				DistDir:          "dist",
				EnableAutoReload: false,
				Config:           map[string]interface{}{"test": "value"},
			}
			fs := NewMockFileSystem()
			if tt.setupFS != nil {
				tt.setupFS(fs)
			}
			renderer := NewMockTemplateRenderer()
			if tt.setupRender != nil {
				tt.setupRender(renderer)
			}
			builder := NewSiteBuilder(site, fs, renderer)

			err := builder.Build()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, fs, renderer)
				}
			}
		})
	}
}

func TestSiteBuilder_buildPages(t *testing.T) {
	tests := []struct {
		name        string
		setupFS     func(*MockFileSystem)
		setupRender func(*MockTemplateRenderer)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Success",
			setupFS:     func(fs *MockFileSystem) {},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: false,
		},
		{
			name: "Mkdir Error",
			setupFS: func(fs *MockFileSystem) {
				fs.mkdirError = errors.New("mkdir error")
			},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: true,
			errorMsg:    "mkdir error",
		},
		{
			name: "Execute Error",
			setupFS: func(fs *MockFileSystem) {},
			setupRender: func(r *MockTemplateRenderer) {
				r.executeError = errors.New("execute error")
			},
			expectError: true,
			errorMsg:    "execute error",
		},
		{
			name: "Read Error",
			setupFS: func(fs *MockFileSystem) {
				fs.readError = errors.New("read error")
			},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: true,
			errorMsg:    "read error",
		},
		{
			name: "Create Error",
			setupFS: func(fs *MockFileSystem) {
				fs.createError = errors.New("create error")
			},
			setupRender: func(r *MockTemplateRenderer) {},
			expectError: true,
			errorMsg:    "create error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := &domain.Site{DistDir: "dist", PagesDir: "pages", Config: map[string]interface{}{"test": "value"}}
			fs := NewMockFileSystem()
			if tt.setupFS != nil {
				tt.setupFS(fs)
			}
			renderer := NewMockTemplateRenderer()
			if tt.setupRender != nil {
				tt.setupRender(renderer)
			}
			builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}
			renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

			err := builder.buildPages(meta.Meta{})

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Check directories created
				expectedMkdirs := []string{"dist", "dist/about", "dist/about/contact"}
				for _, expected := range expectedMkdirs {
					if !contains(fs.mkdirCalls, expected) {
						t.Errorf("Expected mkdir for %s", expected)
					}
				}
			}
		})
	}
}

func TestSiteBuilder_copyAssets(t *testing.T) {
	tests := []struct {
		name        string
		setupFS     func(*MockFileSystem)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Success",
			setupFS:     func(fs *MockFileSystem) {},
			expectError: false,
		},
		{
			name: "Create Error",
			setupFS: func(fs *MockFileSystem) {
				fs.createError = errors.New("create error")
			},
			expectError: true,
			errorMsg:    "create error",
		},
		{
			name: "Walk Error",
			setupFS: func(fs *MockFileSystem) {
				fs.walkError = errors.New("walk error")
			},
			expectError: true,
			errorMsg:    "walk error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := &domain.Site{DistDir: "dist", AssetsDir: "assets"}
			fs := NewMockFileSystem()
			if tt.setupFS != nil {
				tt.setupFS(fs)
			}
			builder := &SiteBuilder{site: site, fs: fs}

			err := builder.copyAssets()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Check mkdir for assets
				expectedDirs := []string{"dist/assets", "dist/assets/css", "dist/assets/js"}
				for _, dir := range expectedDirs {
					if !contains(fs.mkdirCalls, dir) {
						t.Errorf("MkdirAll for %s not called", dir)
					}
				}
			}
		})
	}
}



func TestSiteBuilder_copyFile(t *testing.T) {
	tests := []struct {
		name        string
		setupFS     func(*MockFileSystem)
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupFS: func(fs *MockFileSystem) {
				fs.files["src"] = []byte("test content")
			},
			expectError: false,
		},
		{
			name: "Read Error",
			setupFS: func(fs *MockFileSystem) {
				fs.readError = errors.New("read error")
			},
			expectError: true,
			errorMsg:    "read error",
		},
		{
			name: "Create Error",
			setupFS: func(fs *MockFileSystem) {
				fs.files["src"] = []byte("test content")
				fs.createError = errors.New("create error")
			},
			expectError: true,
			errorMsg:    "create error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := NewMockFileSystem()
			if tt.setupFS != nil {
				tt.setupFS(fs)
			}
			builder := &SiteBuilder{fs: fs}

			err := builder.copyFile("src", "dst")

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !contains(fs.createCalls, "dst") {
					t.Error("Create not called for dst")
				}
			}
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
