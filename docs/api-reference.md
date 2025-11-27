# API Reference

This document provides comprehensive API documentation for the stw-cli Go packages.

## Package Overview

stw-cli is organized into several packages:

- `cmd/stw`: CLI application entry point
- `internal/application`: Application services (site building, serving)
- `internal/domain`: Domain models and business logic
- `internal/infrastructure`: Infrastructure implementations
- `internal/meta`: SEO metadata handling

## cmd/stw

### Main Function

The main entry point for the CLI application.

```go
func main()
```

Initializes Cobra CLI commands and executes the root command.

### Commands

#### buildCmd

```go
var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Build the static site",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

Builds the static site from source files.

#### serveCmd

```go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Build and serve the static site",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

Builds and serves the site with optional auto-reload.

**Flags:**
- `port` (string): Port to serve on (default: "8080")
- `watch` (bool): Enable auto-reload (default: true)

#### initCmd

```go
var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize Wrangler configuration for deployment",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}
```

Initializes Cloudflare Pages deployment configuration.

## internal/application

### SiteBuilder

Handles building the static site.

```go
type SiteBuilder struct {
    site     *domain.Site
    fs       infrastructure.FileSystem
    renderer infrastructure.TemplateRenderer
}
```

#### NewSiteBuilder

```go
func NewSiteBuilder(site *domain.Site, fs infrastructure.FileSystem, renderer infrastructure.TemplateRenderer) *SiteBuilder
```

Creates a new SiteBuilder instance.

#### Build

```go
func (sb *SiteBuilder) Build() error
```

Builds the static site by:
1. Creating the dist directory
2. Loading site metadata
3. Parsing templates
4. Building all pages
5. Copying assets

### SiteServer

Handles serving the static site.

```go
type SiteServer struct {
    site      *domain.Site
    builder   SiteBuilderInterface
    server    HTTPServerInterface
    port      string
    reloadCh  chan struct{}
    clients   map[http.ResponseWriter]bool
    clientsMu sync.Mutex
}
```

#### NewSiteServer

```go
func NewSiteServer(site *domain.Site, builder SiteBuilderInterface, port string) *SiteServer
```

Creates a new SiteServer instance.

#### Serve

```go
func (ss *SiteServer) Serve() error
```

Starts the HTTP server and file watcher.

## internal/domain

### Site

Represents the static site configuration.

```go
type Site struct {
    PagesDir         string
    TemplatesDir     string
    AssetsDir        string
    DistDir          string
    EnableAutoReload bool
    Config           map[string]interface{}
    ConfigPath       string
}
```

**Fields:**
- `PagesDir`: Directory containing HTML pages (default: "pages")
- `TemplatesDir`: Directory containing templates (default: "templates")
- `AssetsDir`: Directory containing static assets (default: "assets")
- `DistDir`: Output directory (default: "dist")
- `EnableAutoReload`: Whether to enable auto-reload in development
- `Config`: Site configuration from config.yaml
- `ConfigPath`: Path to configuration file (default: "config.yaml")

### Page

Represents a web page with content and metadata.

```go
type Page struct {
    Title   string
    Content template.HTML
    Path    string
    IsDev   bool
    Config  map[string]interface{}
    Meta    meta.Meta
}
```

**Fields:**
- `Title`: Page title
- `Content`: Rendered HTML content
- `Path`: Relative path to the page
- `IsDev`: Whether running in development mode
- `Config`: Site configuration
- `Meta`: SEO metadata

## internal/infrastructure

### Interfaces

#### FileSystem

Defines the interface for file operations.

```go
type FileSystem interface {
    WalkDir(root string, fn fs.WalkDirFunc) error
    ReadFile(filename string) ([]byte, error)
    Create(filename string) (io.WriteCloser, error)
    MkdirAll(path string, perm fs.FileMode) error
    RemoveAll(path string) error
}
```

#### TemplateRenderer

Defines the interface for template rendering.

```go
type TemplateRenderer interface {
    ParseFiles(filenames ...string) (*template.Template, error)
    ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}
```

### Implementations

#### OSFileSystem

Implements FileSystem using the os package.

```go
type OSFileSystem struct{}
```

#### GoTemplateRenderer

Implements TemplateRenderer using html/template.

```go
type GoTemplateRenderer struct {
    tmpl *template.Template
}
```

**Custom Functions:**
- `toJson`: Converts interface{} to JSON string

## internal/meta

### Meta

Represents SEO metadata for a page.

```go
type Meta struct {
    Title              string                 `yaml:"title" json:"title"`
    Description        string                 `yaml:"description" json:"description"`
    Canonical          string                 `yaml:"canonical" json:"canonical"`
    Robots             string                 `yaml:"robots" json:"robots"`
    Keywords           string                 `yaml:"keywords" json:"keywords"`
    OgTitle            string                 `yaml:"og_title" json:"og_title"`
    OgDescription      string                 `yaml:"og_description" json:"og_description"`
    OgImage            string                 `yaml:"og_image" json:"og_image"`
    TwitterTitle       string                 `yaml:"twitter_title" json:"twitter_title"`
    TwitterDescription string                 `yaml:"twitter_description" json:"twitter_description"`
    TwitterImage       string                 `yaml:"twitter_image" json:"twitter_image"`
    JsonLd             map[string]interface{} `yaml:"jsonld" json:"jsonld"`
}
```

### Functions

#### Validate

```go
func (m *Meta) Validate(assetsDir string) error
```

Validates meta fields:
- Title ≤ 60 characters
- Description ≤ 160 characters
- Open Graph image under /assets/

#### Merge

```go
func Merge(siteMeta, pageMeta Meta) Meta
```

Merges site-wide defaults with page-specific overrides.

#### ParseFrontMatter

```go
func ParseFrontMatter(content string) (Meta, string, error)
```

Parses YAML or JSON front matter from page content.

#### LoadSiteMeta

```go
func LoadSiteMeta(config map[string]interface{}) Meta
```

Loads site-wide meta configuration from config map.

## Constants

### Domain Constants

```go
const (
    IndexFile          = "index.html"
    BaseTemplate       = "base.html"
    HeaderTemplateFile = "components/header.html"
    FooterTemplateFile = "components/footer.html"
    HeadTemplateFile   = "partials/head.html"
)
```

## Error Handling

stw-cli uses standard Go error handling. Common error conditions:

- Missing configuration file
- Invalid YAML/JSON syntax
- Missing template files
- File system permission errors
- Template parsing errors
- Metadata validation failures

## Dependencies

### External Dependencies

- `github.com/spf13/cobra`: CLI framework
- `gopkg.in/yaml.v3`: YAML parsing
- `github.com/fsnotify/fsnotify`: File watching

### Go Version

Requires Go 1.19 or later.

## Testing

Run tests with:

```bash
go test ./...
```

### Test Coverage

- `internal/meta`: Metadata parsing and validation
- `internal/application`: Site building and serving
- `internal/infrastructure`: File system and template operations

## Examples

### Basic Site Building

```go
config, err := loadConfig()
if err != nil {
    log.Fatal(err)
}

site := &domain.Site{
    PagesDir:     "pages",
    TemplatesDir: "templates",
    AssetsDir:    "assets",
    DistDir:      "dist",
    Config:       config,
}

fs := &infrastructure.OSFileSystem{}
renderer := &infrastructure.GoTemplateRenderer{}
builder := application.NewSiteBuilder(site, fs, renderer)

if err := builder.Build(); err != nil {
    log.Fatal(err)
}
```

### Serving with Auto-reload

```go
server := application.NewSiteServer(site, builder, "8080")
if err := server.Serve(); err != nil {
    log.Fatal(err)
}
```

### Custom Metadata

```go
pageMeta := meta.Meta{
    Title:       "My Page",
    Description: "Page description",
    OgImage:     "/assets/images/og.jpg",
}

if err := pageMeta.Validate("assets"); err != nil {
    log.Fatal(err)
}
```