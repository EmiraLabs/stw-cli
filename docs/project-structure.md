# Project Structure

This document describes the directory structure of the stw-cli project codebase.

## Root Directory

```
stw-cli/
├── cmd/                  # CLI entry point
├── internal/             # Internal packages
├── docs/                 # Documentation
├── config.yaml           # Example site config
├── wrangler.json         # Example deployment config
├── go.mod                # Go module file
├── go.sum                # Go dependencies
└── README.md             # Project README
```

## Core Directories

### cmd/

Contains the main CLI application entry point.

```
cmd/
└── stw/
    ├── main.go           # Main CLI application
    └── main_test.go      # Tests for main
```

### internal/

Contains the internal Go packages that implement the core functionality.

```
internal/
├── application/          # Application layer (site builder, server)
├── domain/               # Domain models (page, site)
├── infrastructure/       # Infrastructure (filesystem, templating)
└── meta/                 # Metadata helpers
```

#### application/
- `sitebuilder.go`: Logic for building the static site
- `siteserver.go`: HTTP server for serving the site

#### domain/
- `page.go`: Page model and parsing
- `site.go`: Site model

#### infrastructure/
- `gotemplaterenderer.go`: Go template rendering
- `osfilesystem.go`: OS filesystem operations
- `interfaces.go`: Infrastructure interfaces

#### meta/
- `helpers.go`: Helper functions for metadata
- `meta_test.go`: Tests

### docs/

Contains all documentation files.

```
docs/
├── api-reference.md      # API docs
├── commands.md           # CLI commands
├── configuration.md      # Config options
├── contributing.md       # Contributing guide
├── deployment-cloudflare.md  # Cloudflare deployment
├── deployment-manual.md  # Manual deployment
├── installation.md       # Installation
├── project-structure.md  # This file
├── quick-start.md        # Quick start guide
├── README.md             # Docs index
├── seo-meta.md           # SEO meta guide
├── site-structure.md     # Site structure (user-facing)
├── templates.md          # Templates guide
└── troubleshooting.md    # Troubleshooting
```

## Configuration Files

### go.mod

Go module definition with dependencies.

### config.yaml

Example site configuration for the demo site.

### wrangler.json

Example Cloudflare Pages deployment configuration.

## Development Workflow

- Code changes in `cmd/` and `internal/`
- Tests alongside source files
- Documentation in `docs/`
- Build with `go build ./cmd/stw`
- Test with `go test ./...`