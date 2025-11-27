# Installation

This guide covers how to install and set up stw-cli on your system.

## Prerequisites

- Go 1.19 or later installed on your system
- Git (for cloning the repository)

## Installation from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/EmiraLabs/stw-cli.git
   cd stw-cli
   ```

2. **Build the binary:**
   ```bash
   go build -o stw ./cmd/stw
   ```

3. **Verify installation:**
   ```bash
   ./stw --help
   ```

## Installation with Go Install

You can also install stw-cli directly using Go:

```bash
go install github.com/EmiraLabs/stw-cli/cmd/stw@latest
```

This will install the `stw` binary to your `$GOPATH/bin` directory. Make sure this directory is in your `PATH`.

## Development Setup

If you want to contribute to stw-cli or run tests:

1. **Clone and build:**
   ```bash
   git clone https://github.com/EmiraLabs/stw-cli.git
   cd stw-cli
   go build -o stw ./cmd/stw
   ```

2. **Run tests:**
   ```bash
   go test ./...
   ```

3. **Run a specific test:**
   ```bash
   go test ./internal/meta
   ```

## System Requirements

- **Operating System:** Linux, macOS, Windows
- **Go Version:** 1.19+
- **Memory:** Minimal (builds are fast and memory-efficient)
- **Disk Space:** ~10MB for the binary and dependencies

## Updating

To update stw-cli:

1. Pull the latest changes:
   ```bash
   git pull origin main
   ```

2. Rebuild:
   ```bash
   go build -o stw ./cmd/stw
   ```

## Uninstalling

To uninstall stw-cli:

1. Remove the binary:
   ```bash
   rm /path/to/stw
   ```

2. If installed via `go install`, remove from GOPATH:
   ```bash
   rm $(go env GOPATH)/bin/stw
   ```

## Troubleshooting

### Build fails with "command not found"

Ensure Go is installed and `go` is in your PATH:

```bash
go version
```

### Permission denied when running

Make sure the binary has execute permissions:

```bash
chmod +x stw
```

### Tests fail

Ensure all dependencies are available:

```bash
go mod tidy
go mod download
```