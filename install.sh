#!/bin/sh
# stw-cli installation script
# 
# Usage:
#   curl -sSL https://raw.githubusercontent.com/EmiraLabs/stw-cli/main/install.sh | sh
#
# Or with custom installation directory:
#   curl -sSL https://raw.githubusercontent.com/EmiraLabs/stw-cli/main/install.sh | INSTALL_DIR=/custom/path sh

set -e

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get latest version from GitHub
LATEST_VERSION=$(curl -s https://api.github.com/repos/EmiraLabs/stw-cli/releases/latest | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
  echo "Error: Could not determine latest version"
  exit 1
fi

echo "Installing stw-cli ${LATEST_VERSION}..."

# Determine installation directory
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Download URL
BINARY_NAME="stw"
if [ "$OS" = "windows" ]; then
  BINARY_NAME="stw.exe"
fi

ARCHIVE_NAME="stw-cli_${LATEST_VERSION#v}_${OS}_${ARCH}.tar.gz"
if [ "$OS" = "windows" ]; then
  ARCHIVE_NAME="stw-cli_${LATEST_VERSION#v}_${OS}_${ARCH}.zip"
fi

DOWNLOAD_URL="https://github.com/EmiraLabs/stw-cli/releases/download/${LATEST_VERSION}/${ARCHIVE_NAME}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
echo "Downloading from $DOWNLOAD_URL..."
if command -v curl > /dev/null 2>&1; then
  curl -sL "$DOWNLOAD_URL" -o archive
elif command -v wget > /dev/null 2>&1; then
  wget -q "$DOWNLOAD_URL" -O archive
else
  echo "Error: Neither curl nor wget found"
  exit 1
fi

# Extract
if [ "$OS" = "windows" ]; then
  unzip -q archive
else
  tar -xzf archive
fi

# Install
echo "Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$BINARY_NAME" "$INSTALL_DIR/"
  chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
  sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

# Cleanup
cd -
rm -rf "$TMP_DIR"

echo "✅ stw-cli ${LATEST_VERSION} installed successfully!"
echo ""
echo "Get started:"
echo "  stw init my-site --wrangler"
echo ""
echo "Documentation: https://github.com/EmiraLabs/stw-cli"
