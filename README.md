# stw-cli

A simple static website generator and server CLI tool written in Go.

## Features

- Build static websites from HTML pages and templates
- Serve the built site locally
- Copy static assets automatically
- Auto-reload browser when files change during development

## Installation

Ensure you have Go installed on your system.

Clone the repository and build the binary:

```bash
git clone https://github.com/EmiraLabs/stw-cli.git
cd stw-cli
go build -o stw ./cmd/stw
```

## Usage

### Build the site

```bash
./stw -build
```

This will generate the site in the `dist` directory.

### Serve the site

```bash
./stw -serve
```

This will build the site and serve it on `http://localhost:8080`. The site will automatically rebuild and reload in the browser when HTML, template, or asset files are changed.

To disable auto-reload:

```bash
./stw serve --watch=false
```

## Project Structure

- `pages/`: Directory containing HTML pages
- `templates/`: Directory containing base template and components
- `assets/`: Directory containing static assets (CSS, JS, images, etc.)
- `dist/`: Output directory for the built site (generated)