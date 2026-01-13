#!/usr/bin/env bash
set -euo pipefail

UPCTL_NAME="upctl"
BASE_URL="https://<org>.github.io/cloud-monitoring-sentinel/dist/latest"

echo "🔧 upctl installer (macOS)"
echo

OS="$(uname -s)"
ARCH="$(uname -m)"

if [[ "$OS" != "Darwin" ]]; then
  echo "❌ This installer supports macOS only"
  exit 1
fi

case "$ARCH" in
  arm64) ARCH="arm64" ;;
  x86_64) ARCH="amd64" ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

BINARY="${UPCTL_NAME}-darwin-${ARCH}"

echo "Detected platform: darwin/${ARCH}"
echo

echo "Where would you like to install upctl?"
echo "  1) ~/bin"
echo "  2) ~/.local/bin"
echo "  3) Current directory"
echo
read -rp "Select [1-3]: " choice

case "$choice" in
  1) INSTALL_DIR="$HOME/bin" ;;
  2) INSTALL_DIR="$HOME/.local/bin" ;;
  3) INSTALL_DIR="$(pwd)" ;;
  *) echo "❌ Invalid selection"; exit 1 ;;
esac

mkdir -p "$INSTALL_DIR"
INSTALL_PATH="${INSTALL_DIR}/${UPCTL_NAME}"

echo
echo "Installing to: $INSTALL_PATH"
echo

curl -fsSL "${BASE_URL}/${BINARY}" -o "$INSTALL_PATH"
chmod +x "$INSTALL_PATH"

echo
echo "✅ upctl installed successfully!"
echo

if [[ "$INSTALL_DIR" == "$(pwd)" ]]; then
  echo "Run it with:"
  echo "  ./upctl bench rightsize --cluster eu"
else
  echo "Run it with:"
  echo "  ${INSTALL_PATH} bench rightsize --cluster eu"
fi

echo
echo "Optional convenience:"
echo "  alias upctl=\"${INSTALL_PATH}\""
echo
