#!/usr/bin/env bash
set -euo pipefail

echo "🔧 upctl installer (macOS)"
echo

# --- platform detection ---
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [[ "$OS" != "darwin" ]]; then
  echo "❌ Unsupported OS: $OS (macOS only)"
  exit 1
fi

case "$ARCH" in
  arm64) BIN_ARCH="arm64" ;;
  x86_64) BIN_ARCH="amd64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

echo "Detected platform: darwin/$BIN_ARCH"

# --- install dir (non-interactive default) ---
INSTALL_DIR="${UPCTL_INSTALL_DIR:-$HOME/.local/bin}"
BIN_NAME="upctl"
BASE_URL="https://benjaminvolodarsky.github.io/cloud-monitoring-sentinel"
BIN_URL="$BASE_URL/dist/latest/${BIN_NAME}-darwin-${BIN_ARCH}"

echo "Installing upctl to $INSTALL_DIR"

mkdir -p "$INSTALL_DIR"

# --- download ---
curl -fsSL "$BIN_URL" -o "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

echo "✔ upctl installed successfully"
echo

# --- PATH hint ---
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "ℹ️  Make sure $INSTALL_DIR is in your PATH"
  echo "   You may need to add this line to your shell config:"
  echo "   export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo
echo "Run: upctl doctor"