package application

import (
	"bytes"
	"net/http"
	"sync"
	"testing"

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

func TestSiteServer_Serve(t *testing.T) {
	site := &domain.Site{DistDir: "dist"}
	builder := &mockSiteBuilder{}
	httpServer := &mockHTTPServer{}
	server := &SiteServer{
		site:    site,
		builder: builder,
		server:  httpServer,
		port:    "8080",
	}

	err := server.Serve()
	if err != httpServer.listenError {
		t.Errorf("Expected error %v, got %v", httpServer.listenError, err)
	}
	if !builder.buildCalled {
		t.Error("Build was not called")
	}
	if !httpServer.listenCalled {
		t.Error("ListenAndServe was not called")
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
