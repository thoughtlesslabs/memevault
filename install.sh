#!/bin/bash
set -e

# Detect OS and Arch
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)     BINARY_OS="linux" ;;
    Darwin)    BINARY_OS="darwin" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)  BINARY_ARCH="amd64" ;;
    arm64)   BINARY_ARCH="arm64" ;;
    aarch64) BINARY_ARCH="arm64" ;;
    *) echo "Unsupported Architecture: $ARCH"; exit 1 ;;
esac

# Construct binary name
BINARY_NAME="memevault_${BINARY_OS}_${BINARY_ARCH}"
GITHUB_REPO="thoughtlesslabs/memevault"
LATEST_RELEASE_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/${BINARY_NAME}"

if [ "$BINARY_OS" == "darwin" ] && [ "$BINARY_ARCH" == "amd64" ]; then
    # Fallback to amd64 for Intel Macs or check if universal
    # Actually we just use the specific one.
    true
fi

# Determine install dir
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
    echo "Installing to $INSTALL_DIR (requires adding to PATH if not present)"
fi

echo "Downloading Memevault ($BINARY_OS/$BINARY_ARCH)..."
curl -fsSL "$LATEST_RELEASE_URL" -o "${INSTALL_DIR}/memevault"
chmod +x "${INSTALL_DIR}/memevault"

echo "Successfully installed Memevault to ${INSTALL_DIR}/memevault"
echo "Run 'memevault --help' to get started!"
