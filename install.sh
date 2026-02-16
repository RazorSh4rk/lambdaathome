#!/usr/bin/env bash
set -euo pipefail

REPO="RazorSh4rk/lambdaathome"
INSTALL_DIR="/usr/local/bin"
SERVICE_NAME="lambdaathome"

# detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected architecture: $ARCH"

# get latest release tag
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST" ]; then
    echo "Failed to fetch latest release"
    exit 1
fi

echo "Latest version: $LATEST"

TARBALL="${SERVICE_NAME}_linux_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

echo "Downloading ${URL}..."

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

curl -fsSL "$URL" -o "${TMP}/${TARBALL}"
tar -xzf "${TMP}/${TARBALL}" -C "$TMP"

install -m 755 "${TMP}/${SERVICE_NAME}" "${INSTALL_DIR}/${SERVICE_NAME}"

echo "Installed ${SERVICE_NAME} ${LATEST} to ${INSTALL_DIR}/${SERVICE_NAME}"
