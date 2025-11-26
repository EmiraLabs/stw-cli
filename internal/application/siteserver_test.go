package application

import (
	"net/http"
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
