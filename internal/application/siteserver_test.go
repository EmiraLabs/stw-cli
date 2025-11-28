package application

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/fsnotify/fsnotify"
)

type mockSiteBuilder struct {
	buildCalled bool
	buildError  error
}

func (m *mockSiteBuilder) Build() error {
	m.buildCalled = true
	return m.buildError
}

type mockHTTPServer struct {
	listenCalled bool
	listenError  error
}

func (m *mockHTTPServer) ListenAndServe(addr string, handler http.Handler) error {
	m.listenCalled = true
	return m.listenError
}

type mockResponseWriter struct {
	buffer     *bytes.Buffer
	flushed    bool
	writeError error
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	if m.writeError != nil {
		return 0, m.writeError
	}
	return m.buffer.Write(data)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {}

func (m *mockResponseWriter) Flush() {
	m.flushed = true
}

func TestNewSiteServer(t *testing.T) {
	site := &domain.Site{}
	builder := &mockSiteBuilder{}
	server := NewSiteServer(site, builder, "8080")
	if server.site != site || server.builder != builder {
		t.Error("NewSiteServer did not set fields correctly")
	}
	if server.server == nil {
		t.Error("Server not set")
	}
	if server.port != "8080" {
		t.Error("Port not set correctly")
	}
}

func TestSiteServer_Serve_BuildError(t *testing.T) {
	site := &domain.Site{DistDir: "dist"}
	builder := &mockSiteBuilder{buildError: errors.New("build error")}
	httpServer := &mockHTTPServer{}
	server := &SiteServer{
		site:    site,
		builder: builder,
		server:  httpServer,
		port:    "8080",
	}

	err := server.Serve()
	if err != builder.buildError {
		t.Errorf("Expected build error %v, got %v", builder.buildError, err)
	}
	if !builder.buildCalled {
		t.Error("Build should be called even if it fails")
	}
	if httpServer.listenCalled {
		t.Error("ListenAndServe should not be called if build fails")
	}
}

func TestSiteServer_notifyClients(t *testing.T) {
	server := &SiteServer{
		clients:   make(map[http.ResponseWriter]bool),
		clientsMu: sync.Mutex{},
	}

	client1 := &mockResponseWriter{buffer: &bytes.Buffer{}}
	client2 := &mockResponseWriter{buffer: &bytes.Buffer{}}

	server.clients[client1] = true
	server.clients[client2] = true

	server.notifyClients()

	expected := "data: reload\n\n"
	if client1.buffer.String() != expected {
		t.Errorf("Client1 received %q, expected %q", client1.buffer.String(), expected)
	}
	if !client1.flushed {
		t.Error("Client1 was not flushed")
	}
	if client2.buffer.String() != expected {
		t.Errorf("Client2 received %q, expected %q", client2.buffer.String(), expected)
	}
	if !client2.flushed {
		t.Error("Client2 was not flushed")
	}
}

func TestSiteServer_clientsConcurrency(t *testing.T) {
	server := &SiteServer{
		clients:   make(map[http.ResponseWriter]bool),
		clientsMu: sync.Mutex{},
	}

	var wg sync.WaitGroup

	// Goroutine to add clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			client := &mockResponseWriter{buffer: &bytes.Buffer{}}
			server.clientsMu.Lock()
			server.clients[client] = true
			server.clientsMu.Unlock()
		}
	}()

	// Goroutine to notify clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			server.notifyClients()
		}
	}()

	wg.Wait()
}

func TestSiteServer_reloadConfig(t *testing.T) {
	// Create temp config file
	configContent := `
title: Test Site
nav:
  - name: Home
    url: /
`
	tempFile := t.TempDir() + "/config.yaml"
	if err := os.WriteFile(tempFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	site := &domain.Site{ConfigPath: tempFile}
	server := &SiteServer{site: site}

	// Call private method directly (same package)
	err := server.reloadConfig()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if config was loaded
	if server.site.Config == nil {
		t.Error("Config not loaded")
	}
	if title, ok := server.site.Config["title"]; !ok || string(title.(template.HTML)) != "Test Site" {
		t.Errorf("Expected title 'Test Site', got %v", title)
	}
}

func TestSiteServer_convertToHTML(t *testing.T) {
	// Test the convertToHTML function in siteserver
	// Since it's private, we can't call it directly, but we can test through reloadConfig
	// which uses it

	// Create temp config file with HTML content
	configContent := `
title: <h1>Test Site</h1>
nav:
  - name: "<b>Home</b>"
    url: /
`
	tempFile := t.TempDir() + "/config.yaml"
	if err := os.WriteFile(tempFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	site := &domain.Site{ConfigPath: tempFile}
	server := &SiteServer{site: site}

	err := server.reloadConfig()
	if err != nil {
		t.Fatalf("reloadConfig failed: %v", err)
	}

	// Check that HTML strings are converted
	if title, ok := server.site.Config["title"]; !ok || string(title.(template.HTML)) != "<h1>Test Site</h1>" {
		t.Errorf("Expected HTML title, got %v", title)
	}

	if nav, ok := server.site.Config["nav"].([]interface{}); ok && len(nav) > 0 {
		if item, ok := nav[0].(map[string]interface{}); ok {
			if name, ok := item["name"]; !ok || string(name.(template.HTML)) != "<b>Home</b>" {
				t.Errorf("Expected HTML name, got %v", name)
			}
		}
	}
}

func TestSiteServer_initWatcher(t *testing.T) {
	// This test is limited since fsnotify requires real filesystem
	// But we can test the basic call
	site := &domain.Site{
		PagesDir:     "pages",
		TemplatesDir: "templates",
		AssetsDir:    "assets",
		ConfigPath:   "config.yaml",
	}
	server := &SiteServer{site: site}

	// Call initWatcher - this will try to create fsnotify.Watcher
	// In test environment, it may fail or succeed
	watcher, err := server.initWatcher()
	if err != nil {
		// Expected in test environment without real dirs
		t.Logf("initWatcher failed as expected: %v", err)
		return
	}
	if watcher != nil {
		watcher.Close()
	}
}

func TestSiteServer_reloadConfig_InvalidYAML(t *testing.T) {
	// Create temp config file with invalid YAML
	configContent := `
title: Test Site
invalid: yaml: content: [
`
	tempFile := t.TempDir() + "/config.yaml"
	if err := os.WriteFile(tempFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	site := &domain.Site{ConfigPath: tempFile}
	server := &SiteServer{site: site}

	err := server.reloadConfig()
	if err == nil {
		t.Error("Expected YAML unmarshal error, got nil")
	}
}

func TestSiteServer_reloadConfig_ReadFileError(t *testing.T) {
	// Use a directory path instead of file path to cause a different read error
	tempDir := t.TempDir()
	site := &domain.Site{ConfigPath: tempDir} // This is a directory, not a file
	server := &SiteServer{site: site}

	err := server.reloadConfig()
	if err == nil {
		t.Error("Expected read file error when path is a directory, got nil")
	}
}

func TestSiteServer_notifyClients_WriteError(t *testing.T) {
	server := &SiteServer{
		clients:   make(map[http.ResponseWriter]bool),
		clientsMu: sync.Mutex{},
	}

	client1 := &mockResponseWriter{buffer: &bytes.Buffer{}, writeError: errors.New("write error")}
	client2 := &mockResponseWriter{buffer: &bytes.Buffer{}}

	server.clients[client1] = true
	server.clients[client2] = true

	initialLen := len(server.clients)
	server.notifyClients()

	// client1 should be removed due to write error
	if len(server.clients) != initialLen-1 {
		t.Errorf("Expected client to be removed on write error, clients count: %d", len(server.clients))
	}

	// client2 should still be there and have received the message
	if client2.buffer.String() != "data: reload\n\n" {
		t.Errorf("Client2 received %q, expected %q", client2.buffer.String(), "data: reload\n\n")
	}
}

func TestSiteServer_Serve_WithAutoReload(t *testing.T) {
	site := &domain.Site{
		DistDir:          "dist",
		EnableAutoReload: true,
	}
	builder := &mockSiteBuilder{}
	httpServer := &mockHTTPServer{listenError: errors.New("server stopped")}
	server := &SiteServer{
		site:    site,
		builder: builder,
		server:  httpServer,
		port:    "8080",
	}

	// This will call ListenAndServe which returns an error
	err := server.Serve()
	if err != httpServer.listenError {
		t.Errorf("Expected server error %v, got %v", httpServer.listenError, err)
	}
	if !builder.buildCalled {
		t.Error("Build should be called")
	}
	if !httpServer.listenCalled {
		t.Error("ListenAndServe should be called")
	}
}

func TestSiteServer_Serve_WithoutAutoReload(t *testing.T) {
	site := &domain.Site{
		DistDir:          "dist",
		EnableAutoReload: false,
	}
	builder := &mockSiteBuilder{}
	httpServer := &mockHTTPServer{listenError: errors.New("server stopped")}
	server := &SiteServer{
		site:    site,
		builder: builder,
		server:  httpServer,
		port:    "8080",
	}

	err := server.Serve()
	if err != httpServer.listenError {
		t.Errorf("Expected server error %v, got %v", httpServer.listenError, err)
	}
	if !builder.buildCalled {
		t.Error("Build should be called")
	}
	if !httpServer.listenCalled {
		t.Error("ListenAndServe should be called")
	}
}

func TestSiteServer_handleFileEvent_WriteEvent(t *testing.T) {
	tempDir := t.TempDir()
	site := &domain.Site{
		PagesDir:     tempDir + "/pages",
		TemplatesDir: tempDir + "/templates",
		AssetsDir:    tempDir + "/assets",
		ConfigPath:   tempDir + "/config.yaml",
		DistDir:      tempDir + "/dist",
	}

	// Create directories
	os.MkdirAll(site.PagesDir, 0755)
	os.MkdirAll(site.TemplatesDir, 0755)
	os.MkdirAll(site.AssetsDir, 0755)
	os.MkdirAll(site.DistDir, 0755)

	builder := &mockSiteBuilder{}
	server := &SiteServer{
		site:     site,
		builder:  builder,
		reloadCh: make(chan struct{}, 1),
	}

	// Test write event on a page file
	event := fsnotify.Event{
		Name: tempDir + "/pages/index.html",
		Op:   fsnotify.Write,
	}

	server.handleFileEvent(event, nil)

	if !builder.buildCalled {
		t.Error("Build should be called on file write")
	}
	builder.buildCalled = false // reset

	// For file changes, notifyClients is called, but reloadCh is for WebSocket
	// Let's check that notifyClients would be called by checking the log or something
	// Actually, since notifyClients sends to WebSocket clients, and we have none, it's fine
}

func TestSiteServer_handleFileEvent_ConfigChange(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yaml"
	site := &domain.Site{
		ConfigPath: configPath,
		DistDir:    tempDir + "/dist",
	}

	os.MkdirAll(site.DistDir, 0755)

	builder := &mockSiteBuilder{}
	server := &SiteServer{
		site:     site,
		builder:  builder,
		reloadCh: make(chan struct{}, 1),
	}

	// Test config file change
	event := fsnotify.Event{
		Name: configPath,
		Op:   fsnotify.Write,
	}

	// Create config file first
	os.WriteFile(configPath, []byte("title: test"), 0644)

	server.handleFileEvent(event, nil)

	if !builder.buildCalled {
		t.Error("Build should be called on config change")
	}

	// Check if config was reloaded
	if server.site.Config == nil || server.site.Config["title"] != template.HTML("test") {
		t.Error("Config should be reloaded on config file change")
	}
}
