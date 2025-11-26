package application

import (
	"log"
	"net/http"

	"github.com/EmiraLabs/stw-cli/internal/domain"
)

// SiteServer handles serving the static site
type SiteServer struct {
	site    *domain.Site
	builder *SiteBuilder
}

// NewSiteServer creates a new SiteServer
func NewSiteServer(site *domain.Site, builder *SiteBuilder) *SiteServer {
	return &SiteServer{
		site:    site,
		builder: builder,
	}
}

// Serve builds and serves the site
func (ss *SiteServer) Serve() error {
	if err := ss.builder.Build(); err != nil {
		return err
	}
	fs := http.FileServer(http.Dir(ss.site.DistDir))
	log.Printf("Serving %s on http://localhost:8001", ss.site.DistDir)
	return http.ListenAndServe(":8001", fs)
}
