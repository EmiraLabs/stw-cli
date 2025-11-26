package application

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/EmiraLabs/stw-cli/internal/domain"
	"github.com/fsnotify/fsnotify"
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
	site      *domain.Site
	builder   SiteBuilderInterface
	server    HTTPServerInterface
	port      string
	reloadCh  chan struct{}
	clients   map[http.ResponseWriter]bool
	clientsMu sync.Mutex
}

// NewSiteServer creates a new SiteServer
func NewSiteServer(site *domain.Site, builder SiteBuilderInterface, port string) *SiteServer {
	return &SiteServer{
		site:     site,
		builder:  builder,
		server:   &DefaultHTTPServer{},
		port:     port,
		reloadCh: make(chan struct{}, 1),
		clients:  make(map[http.ResponseWriter]bool),
	}
}

// Serve builds and serves the site
func (ss *SiteServer) Serve() error {
	if err := ss.builder.Build(); err != nil {
		return err
	}

	// Start file watcher if enabled
	if ss.site.EnableAutoReload {
		go ss.watchFiles()
	}

	// Custom handler
	mux := http.NewServeMux()
	if ss.site.EnableAutoReload {
		mux.HandleFunc("/__reload", ss.handleReload)
	}
	mux.Handle("/", http.FileServer(http.Dir(ss.site.DistDir)))

	log.Printf("Serving %s on http://localhost:%s", ss.site.DistDir, ss.port)
	return ss.server.ListenAndServe(":"+ss.port, mux)
}

func (ss *SiteServer) handleReload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Client connected to /__reload")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ss.clientsMu.Lock()
	ss.clients[w] = true
	ss.clientsMu.Unlock()

	// Remove client on disconnect
	defer func() {
		ss.clientsMu.Lock()
		delete(ss.clients, w)
		ss.clientsMu.Unlock()
	}()

	// Listen for reload events
	for {
		select {
		case <-ss.reloadCh:
			log.Printf("Sending reload event to client")
			if _, err := w.Write([]byte("data: reload\n\n")); err != nil {
				return
			}
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (ss *SiteServer) watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	addDir := func(path string) error {
		return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return watcher.Add(p)
			}
			return nil
		})
	}

	dirs := []string{ss.site.PagesDir, ss.site.TemplatesDir, ss.site.AssetsDir}
	for _, dir := range dirs {
		if err := addDir(dir); err != nil {
			log.Printf("Error watching %s: %v", dir, err)
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				// If a directory is created, add it to watch
				if event.Has(fsnotify.Create) {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						watcher.Add(event.Name)
					}
				}
				log.Printf("File changed: %s", event.Name)
				if err := ss.builder.Build(); err != nil {
					log.Printf("Build error: %v", err)
				} else {
					// Notify clients
					ss.clientsMu.Lock()
					log.Printf("Notifying %d clients", len(ss.clients))
					for client := range ss.clients {
						if _, err := client.Write([]byte("data: reload\n\n")); err != nil {
							delete(ss.clients, client)
						} else {
							client.(http.Flusher).Flush()
						}
					}
					ss.clientsMu.Unlock()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
