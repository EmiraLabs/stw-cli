package application

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/EmiraLabs/stw-cli/internal/domain"
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

func TestSiteServer_handleReload(t *testing.T) {
	site := &domain.Site{}
	server := &SiteServer{
		site:    site,
		clients: make(map[http.ResponseWriter]bool),
	}

	// Create a test request with cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/__reload", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Call handleReload in a goroutine
	go server.handleReload(w, req)

	// Wait a bit, then cancel the request
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Wait a bit for the handler to finish
	time.Sleep(10 * time.Millisecond)

	// Check response headers
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Error("Wrong content type")
	}
}
