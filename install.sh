#!/bin/bash

set -e

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Map to GoReleaser naming
case "$OS" in
    Darwin)
        OS_NAME="Darwin"
        ;;
    Linux)
        OS_NAME="Linux"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCH_NAME="x86_64"
        ;;
    aarch64|arm64)
        ARCH_NAME="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

REPO="NoaTamburrini/portman"
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="${INSTALL_DIR}/portman"

echo "Installing portman for ${OS_NAME}/${ARCH_NAME}..."

# Get latest release version
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi

echo "Latest version: ${LATEST_VERSION}"

# Archive name from GoReleaser
ARCHIVE_NAME="portman_${OS_NAME}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${ARCHIVE_NAME}"

echo "Downloading from: ${DOWNLOAD_URL}"

# Download and extract
TMP_DIR=$(mktemp -d)
trap "rm -rf ${TMP_DIR}" EXIT

if ! curl -L -o "${TMP_DIR}/portman.tar.gz" "${DOWNLOAD_URL}"; then
    echo "Failed to download portman"
    exit 1
fi

# Extract
tar -xzf "${TMP_DIR}/portman.tar.gz" -C "${TMP_DIR}"

# Install (may need sudo)
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_DIR}/portman" "${INSTALL_PATH}"
    chmod +x "${INSTALL_PATH}"
else
    echo "Installing to ${INSTALL_PATH} (requires sudo)..."
    sudo mv "${TMP_DIR}/portman" "${INSTALL_PATH}"
    sudo chmod +x "${INSTALL_PATH}"
fi

echo "✅ portman installed successfully to ${INSTALL_PATH}"
echo "Run 'portman' to get started!"
