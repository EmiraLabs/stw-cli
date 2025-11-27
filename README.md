# stw-cli

A simple static website generator and server CLI tool written in Go.

## Features

- Build static websites from HTML pages and templates
- Serve the built site locally with auto-reload
- Copy static assets automatically
- Deploy to Cloudflare Pages using Wrangler
- SEO Meta Support with front matter

## Documentation

Complete documentation is available in the [docs/](docs/) directory:

- [Installation](docs/installation.md)
- [Quick Start](docs/quick-start.md)
- [Configuration](docs/configuration.md)
- [Commands](docs/commands.md)
- [Project Structure](docs/project-structure.md)
- [Templates](docs/templates.md)
- [SEO Meta System](docs/seo-meta.md)
- [Cloudflare Pages Deployment](docs/deployment-cloudflare.md)
- [Manual Deployment](docs/deployment-manual.md)
- [API Reference](docs/api-reference.md)
- [Contributing](docs/contributing.md)
- [Troubleshooting](docs/troubleshooting.md)

## Installation

Ensure you have Go installed on your system.

```bash
git clone https://github.com/EmiraLabs/stw-cli.git
cd stw-cli
go build -o stw ./cmd/stw
```

## Usage

### Build the site

```bash
./stw build
```

### Serve the site

```bash
./stw serve
```

This serves the site on `http://localhost:8080` with auto-reload.

For more details, see the [Quick Start](docs/quick-start.md) and [Commands](docs/commands.md).