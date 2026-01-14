#!/bin/bash
set -e

INSTALL_DIR="$HOME/.envault/bin"
mkdir -p "$INSTALL_DIR"

echo "Installing Envault to $INSTALL_DIR..."

# In production, detect OS/Arch and download
# OS="$(uname -s)"
# ARCH="$(uname -m)"
# DOWNLOAD_URL="..."
# curl -L "$DOWNLOAD_URL" -o "$INSTALL_DIR/envault"

# Local Dev Fallback
if [ -f "./envault" ]; then
    cp ./envault "$INSTALL_DIR/envault"
    echo "Copied local binary."
else
    echo "This is a placeholder installer. In production, I would download the binary."
    # Fail for now if no binary
    exit 1
fi

chmod +x "$INSTALL_DIR/envault"

# Add to PATH
SHELL_CONFIG="$HOME/.bashrc"
if [ -n "$ZSH_VERSION" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Adding to PATH in $SHELL_CONFIG..."
    echo "" >> "$SHELL_CONFIG"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_CONFIG"
    echo "Added to PATH. Run 'source $SHELL_CONFIG' or restart shell."
else
    echo "Already in PATH."
fi

echo "Installation Complete!"
