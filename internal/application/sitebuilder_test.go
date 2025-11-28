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

func TestSiteBuilder_Build_ParseError(t *testing.T) {
	site := &domain.Site{
		PagesDir:     "pages",
		TemplatesDir: "templates",
		AssetsDir:    "assets",
		DistDir:      "dist",
	}
	fs := NewMockFileSystem()
	renderer := NewMockTemplateRenderer()
	renderer.parseError = errors.New("parse error")
	builder := NewSiteBuilder(site, fs, renderer)

	err := builder.Build()
	if err == nil {
		t.Errorf("Expected error from ParseFiles, but got: %v", err)
	}
}

func TestSiteBuilder_Build_MkdirError(t *testing.T) {
	site := &domain.Site{
		PagesDir:     "pages",
		TemplatesDir: "templates",
		AssetsDir:    "assets",
		DistDir:      "dist",
	}
	fs := NewMockFileSystem()
	fs.mkdirError = errors.New("mkdir error")
	renderer := NewMockTemplateRenderer()
	builder := NewSiteBuilder(site, fs, renderer)

	err := builder.Build()
	if err == nil {
		t.Errorf("Expected error from MkdirAll, but got: %v", err)
	}
}

func TestSiteBuilder_Build_RemoveError(t *testing.T) {
	site := &domain.Site{
		PagesDir:     "pages",
		TemplatesDir: "templates",
		AssetsDir:    "assets",
		DistDir:      "dist",
	}
	fs := NewMockFileSystem()
	fs.removeError = errors.New("remove error")
	renderer := NewMockTemplateRenderer()
	builder := NewSiteBuilder(site, fs, renderer)

	err := builder.Build()
	if err == nil {
		t.Errorf("Expected error from RemoveAll, but got: %v", err)
	}
}

func TestSiteBuilder_Build(t *testing.T) {
	site := &domain.Site{
		PagesDir:         "pages",
		TemplatesDir:     "templates",
		AssetsDir:        "assets",
		DistDir:          "dist",
		EnableAutoReload: false,
		Config:           map[string]interface{}{"test": "value"},
	}
	fs := NewMockFileSystem()
	renderer := NewMockTemplateRenderer()
	builder := NewSiteBuilder(site, fs, renderer)

	err := builder.Build()
	if err != nil {
		t.Errorf("Build failed: %v", err)
	}

	// Check that RemoveAll was called
	if len(fs.removeCalls) != 1 || fs.removeCalls[0] != "dist" {
		t.Error("RemoveAll not called correctly")
	}

	// Check MkdirAll for dist
	if len(fs.mkdirCalls) < 1 || fs.mkdirCalls[0] != "dist" {
		t.Error("MkdirAll for dist not called")
	}

	// Check ParseFiles
	if len(renderer.parseFilesCalls) != 1 {
		t.Error("ParseFiles not called")
	}

	// Check that pages were built
	// Should have created dist/index.html, dist/about/index.html, etc.
	expectedCreates := []string{"dist/index.html", "dist/about/index.html", "dist/about/contact/index.html"}
	for _, expected := range expectedCreates {
		found := false
		for _, call := range fs.createCalls {
			if call == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected create call for %s", expected)
		}
	}

	// Check assets copied
	expectedAssetCreates := []string{"dist/assets/css/style.css", "dist/assets/js/app.js"}
	for _, expected := range expectedAssetCreates {
		found := false
		for _, call := range fs.createCalls {
			if call == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected create call for %s", expected)
		}
	}
}

func TestSiteBuilder_buildPages(t *testing.T) {
	site := &domain.Site{DistDir: "dist", PagesDir: "pages", Config: map[string]interface{}{"test": "value"}}
	fs := NewMockFileSystem()
	renderer := NewMockTemplateRenderer()
	builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}
	tmpl, _ := renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

	err := builder.buildPages(tmpl, meta.Meta{})
	if err != nil {
		t.Errorf("buildPages failed: %v", err)
	}

	// Check directories created
	expectedMkdirs := []string{"dist", "dist/about", "dist/about/contact"}
	for _, expected := range expectedMkdirs {
		found := false
		for _, call := range fs.mkdirCalls {
			if call == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected mkdir for %s", expected)
		}
	}
}

func TestSiteBuilder_buildPages_MkdirError(t *testing.T) {
	site := &domain.Site{DistDir: "dist", PagesDir: "pages", Config: map[string]interface{}{}}
	fs := NewMockFileSystem()
	fs.mkdirError = errors.New("mkdir error")
	renderer := NewMockTemplateRenderer()
	builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}
	tmpl, _ := renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

	err := builder.buildPages(tmpl, meta.Meta{})
	if err == nil {
		t.Errorf("Expected error from MkdirAll, but got: %v", err)
	}
}

func TestSiteBuilder_buildPages_ExecuteError(t *testing.T) {
	site := &domain.Site{DistDir: "dist", PagesDir: "pages", Config: map[string]interface{}{}}
	fs := NewMockFileSystem()
	renderer := NewMockTemplateRenderer()
	renderer.executeError = errors.New("execute error")
	builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}
	tmpl, _ := renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

	err := builder.buildPages(tmpl, meta.Meta{})
	if err == nil {
		t.Errorf("Expected error from ExecuteTemplate, but got: %v", err)
	}
}

func TestSiteBuilder_buildPages_ReadError(t *testing.T) {
	site := &domain.Site{DistDir: "dist", PagesDir: "pages", Config: map[string]interface{}{}}
	fs := NewMockFileSystem()
	fs.readError = errors.New("read error")
	renderer := NewMockTemplateRenderer()
	builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}
	tmpl, _ := renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

	err := builder.buildPages(tmpl, meta.Meta{})
	if err == nil {
		t.Errorf("Expected error from ReadFile, but got: %v", err)
	}
}

func TestSiteBuilder_buildPages_CreateError(t *testing.T) {
	site := &domain.Site{DistDir: "dist", PagesDir: "pages", Config: map[string]interface{}{}}
	fs := NewMockFileSystem()
	fs.createError = errors.New("create error")
	renderer := NewMockTemplateRenderer()
	builder := &SiteBuilder{site: site, fs: fs, renderer: renderer}
	tmpl, _ := renderer.ParseFiles(filepath.Join("templates", domain.BaseTemplate))

	err := builder.buildPages(tmpl, meta.Meta{})
	if err == nil {
		t.Errorf("Expected error from Create, but got: %v", err)
	}
}

func TestSiteBuilder_copyAssets(t *testing.T) {
	site := &domain.Site{DistDir: "dist", AssetsDir: "assets"}
	fs := NewMockFileSystem()
	builder := &SiteBuilder{site: site, fs: fs}

	err := builder.copyAssets()
	if err != nil {
		t.Errorf("copyAssets failed: %v", err)
	}

	// Check mkdir for assets
	if !contains(fs.mkdirCalls, "dist/assets") {
		t.Error("MkdirAll for dist/assets not called")
	}
	if !contains(fs.mkdirCalls, "dist/assets/css") {
		t.Error("MkdirAll for dist/assets/css not called")
	}
	if !contains(fs.mkdirCalls, "dist/assets/js") {
		t.Error("MkdirAll for dist/assets/js not called")
	}
}

func TestSiteBuilder_copyAssets_Error(t *testing.T) {
	site := &domain.Site{DistDir: "dist", AssetsDir: "assets"}
	fs := NewMockFileSystem()
	fs.createError = errors.New("create error")
	builder := &SiteBuilder{site: site, fs: fs}

	err := builder.copyAssets()
	if err == nil {
		t.Errorf("Expected error from copyFile, but got: %v", err)
	}
}

func TestSiteBuilder_copyAssets_WalkError(t *testing.T) {
	site := &domain.Site{DistDir: "dist", AssetsDir: "assets"}
	fs := NewMockFileSystem()
	fs.walkError = errors.New("walk error")
	builder := &SiteBuilder{site: site, fs: fs}

	err := builder.copyAssets()
	if err == nil {
		t.Errorf("Expected error from WalkDir, but got: %v", err)
	}
}

func TestSiteBuilder_Build_NoPages(t *testing.T) {
	site := &domain.Site{
		PagesDir:     "pages",
		TemplatesDir: "templates",
		AssetsDir:    "assets",
		DistDir:      "dist",
	}
	fs := NewMockFileSystem()
	fs.hasPages = false
	renderer := NewMockTemplateRenderer()
	builder := NewSiteBuilder(site, fs, renderer)

	err := builder.Build()
	if err != nil {
		t.Errorf("Build failed: %v", err)
	}

	// Check that no pages were built
	// Should not have created any dist files from pages
	for _, call := range fs.createCalls {
		if strings.HasPrefix(call, "dist/") && strings.HasSuffix(call, ".html") {
			t.Errorf("Unexpected create call for page: %s", call)
		}
	}
}

func TestSiteBuilder_copyFile(t *testing.T) {
	fs := NewMockFileSystem()
	fs.files["src"] = []byte("test content")
	builder := &SiteBuilder{fs: fs}

	err := builder.copyFile("src", "dst")
	if err != nil {
		t.Errorf("copyFile failed: %v", err)
	}

	if !contains(fs.createCalls, "dst") {
		t.Error("Create not called for dst")
	}
}

func TestSiteBuilder_copyFile_ReadError(t *testing.T) {
	fs := NewMockFileSystem()
	fs.readError = errors.New("read error")
	builder := &SiteBuilder{fs: fs}

	err := builder.copyFile("src", "dst")
	if err == nil {
		t.Errorf("Expected error from ReadFile, but got: %v", err)
	}
}

func TestSiteBuilder_copyFile_CreateError(t *testing.T) {
	fs := NewMockFileSystem()
	fs.files["src"] = []byte("test content")
	fs.createError = errors.New("create error")
	builder := &SiteBuilder{fs: fs}

	err := builder.copyFile("src", "dst")
	if err == nil {
		t.Errorf("Expected error from Create, but got: %v", err)
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
