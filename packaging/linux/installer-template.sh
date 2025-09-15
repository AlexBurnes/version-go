#!/bin/sh
set -eu

# Version and platform info - will be replaced by build script
VERSION="{{VERSION}}"
PLATFORM="{{PLATFORM}}"
ARCH="{{ARCH}}"
BINARY_NAME="{{BINARY_NAME}}"

# Default install directory
# When piped to sh, arguments don't work the same way
# So we'll use environment variable or default
if test -n "${INSTALL_DIR:-}"; then
    # Use environment variable if set
    INSTALL_DIR="$INSTALL_DIR"
elif test -n "${1:-}"; then
    # Use first argument if available
    INSTALL_DIR="$1"
else
    # Default directory
    INSTALL_DIR="/usr/local/bin"
fi

# GitHub release URL
RELEASE_URL="https://github.com/AlexBurnes/version-go/releases/download/v${VERSION}/version_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"

echo "[*] Installing version CLI v${VERSION} for ${PLATFORM}-${ARCH}"
echo "[*] Target directory: ${INSTALL_DIR}"

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf '$TEMP_DIR'" EXIT

echo "[*] Downloading archive from GitHub..."
if ! wget -q -O "${TEMP_DIR}/archive.tar.gz" "$RELEASE_URL"; then
    echo "[ERROR] Failed to download archive from: $RELEASE_URL"
    echo "[ERROR] Please check if the release exists and try again"
    exit 1
fi

echo "[*] Extracting archive..."
cd "$TEMP_DIR"
if ! tar -xzf archive.tar.gz; then
    echo "[ERROR] Failed to extract archive"
    exit 1
fi

# Check if install.sh exists in the archive
if test -f "install.sh"; then
    echo "[*] Using install.sh from archive"
    chmod +x install.sh
    ./install.sh "$INSTALL_DIR"
else
    echo "[ERROR] install.sh not found in archive"
    echo "[ERROR] Contents of archive:"
    find . -type f | head -10
    exit 1
fi