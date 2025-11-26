package application

import (
	"log"
	"net/http"

	"github.com/EmiraLabs/stw-cli/internal/domain"
)

type SiteBuilderInterface interface {
	Build() error
}

type HTTPServerInterface interface {
	ListenAndServe(addr string, handler http.Handler) error
}

type DefaultHTTPServer struct{}

func (d *DefaultHTTPServer) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

// SiteServer handles serving the static site
type SiteServer struct {
	site    *domain.Site
	builder SiteBuilderInterface
	server  HTTPServerInterface
}

// NewSiteServer creates a new SiteServer
func NewSiteServer(site *domain.Site, builder SiteBuilderInterface) *SiteServer {
	return &SiteServer{
		site:    site,
		builder: builder,
		server:  &DefaultHTTPServer{},
	}
}

// Serve builds and serves the site
func (ss *SiteServer) Serve() error {
	if err := ss.builder.Build(); err != nil {
		return err
	}
	fs := http.FileServer(http.Dir(ss.site.DistDir))
	log.Printf("Serving %s on http://localhost:8001", ss.site.DistDir)
	return ss.server.ListenAndServe(":8001", fs)
}
