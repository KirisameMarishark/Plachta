#!/usr/bin/env bash
set -euo pipefail

REPO="KirisameMarishark/Plachta"
VERSION="v0.1.0"
INSTALL_PATH="/usr/local/bin/plachta"

log_info() {
    echo "[INFO] $*"
}

log_success() {
    echo "[OK] $*"
}

log_error() {
    echo "[ERROR] $*" >&2
}

if [[ "$(uname -s)" != "Linux" ]]; then
    log_error "Plachta binary installer currently supports Linux only."
    exit 1
fi

ARCH="$(uname -m)"

case "$ARCH" in
    x86_64)
        ASSET="plachta-linux-amd64"
        ;;
    aarch64|arm64)
        ASSET="plachta-linux-arm64"
        ;;
    *)
        log_error "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"

log_info "Installing Plachta ${VERSION}"
log_info "Architecture: ${ARCH}"
log_info "Downloading: ${ASSET}"

TMP_FILE="$(mktemp)"

cleanup() {
    rm -f "$TMP_FILE"
}

trap cleanup EXIT

if command -v curl >/dev/null 2>&1; then
    curl -fL --retry 3 --retry-delay 1 \
        -o "$TMP_FILE" \
        "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$TMP_FILE" "$DOWNLOAD_URL"
else
    log_error "Neither curl nor wget is installed."
    exit 1
fi

if [[ ! -s "$TMP_FILE" ]]; then
    log_error "Downloaded file is empty."
    exit 1
fi

chmod +x "$TMP_FILE"

log_info "Installing to ${INSTALL_PATH}"

if [[ ! -w "$(dirname "$INSTALL_PATH")" ]]; then
    if [[ $EUID -eq 0 ]]; then
        mv "$TMP_FILE" "$INSTALL_PATH"
    else
        sudo mv "$TMP_FILE" "$INSTALL_PATH"
    fi
else
    mv "$TMP_FILE" "$INSTALL_PATH"
fi

chmod 0755 "$INSTALL_PATH"

log_success "Plachta installed successfully."

"$INSTALL_PATH" version
