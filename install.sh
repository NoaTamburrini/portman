#!/bin/bash

set -e

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture names
case "$ARCH" in
    x86_64)
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

# Check supported OS
if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
    echo "Unsupported OS: $OS"
    exit 1
fi

REPO="NoaTamburrini/portman"
BINARY_NAME="portman-${OS}-${ARCH}"
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="${INSTALL_DIR}/portman"

echo "Installing portman for ${OS}/${ARCH}..."

# Get latest release version
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi

echo "Latest version: ${LATEST_VERSION}"

# Download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"

echo "Downloading from: ${DOWNLOAD_URL}"

# Download to temp file
TMP_FILE=$(mktemp)
trap "rm -f ${TMP_FILE}" EXIT

if ! curl -L -o "${TMP_FILE}" "${DOWNLOAD_URL}"; then
    echo "Failed to download portman"
    exit 1
fi

# Install (may need sudo)
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_FILE}" "${INSTALL_PATH}"
    chmod +x "${INSTALL_PATH}"
else
    echo "Installing to ${INSTALL_PATH} (requires sudo)..."
    sudo mv "${TMP_FILE}" "${INSTALL_PATH}"
    sudo chmod +x "${INSTALL_PATH}"
fi

echo "✅ portman installed successfully to ${INSTALL_PATH}"
echo "Run 'portman' to get started!"
