#!/usr/bin/env bash
set -euo pipefail

REPO="KirisameMarishark/Plachta"
VERSION="${1:-latest}"
INSTALL_DIR="/opt/plachta"
BIN_LINK="/usr/local/bin/plachta"

log_info() {
    printf '\033[0;32m[INFO]\033[0m %s\n' "$1"
}

log_error() {
    printf '\033[0;31m[ERROR]\033[0m %s\n' "$1" >&2
}

if [[ "$(uname -s)" != "Linux" ]]; then
    log_error "This installer supports Linux only."
    exit 1
fi

ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ASSET="plachta-linux-amd64.tar.gz" ;;
    aarch64|arm64) ASSET="plachta-linux-arm64.tar.gz" ;;
    *)
        log_error "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

if [[ "$VERSION" == "latest" ]]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

if [[ "$EUID" -ne 0 ]]; then
    log_error "Please run as root: sudo bash install.sh [version]"
    exit 1
fi

for cmd in curl tar install; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        log_error "Required command not found: $cmd"
        exit 1
    fi
done

TMP_DIR="$(mktemp -d)"
STAGE_DIR=""
BACKUP_DIR=""
DEPLOYED=0

cleanup() {
    if [[ "$DEPLOYED" -eq 0 && -n "$BACKUP_DIR" && -d "$BACKUP_DIR" && ! -e "$INSTALL_DIR" ]]; then
        mv "$BACKUP_DIR" "$INSTALL_DIR" || true
    fi
    if [[ -n "$STAGE_DIR" && -d "$STAGE_DIR" ]]; then
        rm -rf "$STAGE_DIR"
    fi
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

log_info "Downloading Plachta package: ${VERSION}"
if ! curl -fL --retry 3 "$DOWNLOAD_URL" -o "$TMP_DIR/plachta.tar.gz"; then
    log_error "Download failed: $DOWNLOAD_URL"
    exit 1
fi

if ! tar -xzf "$TMP_DIR/plachta.tar.gz" -C "$TMP_DIR"; then
    log_error "Could not extract release package."
    exit 1
fi

if [[ ! -f "$TMP_DIR/plachta" || ! -f "$TMP_DIR/VERSION" || ! -d "$TMP_DIR/modules" ]]; then
    log_error "Invalid release package: plachta, VERSION, or modules/ is missing."
    exit 1
fi

PACKAGE_VERSION="$(tr -d '[:space:]' < "$TMP_DIR/VERSION")"
if [[ -z "$PACKAGE_VERSION" ]]; then
    log_error "VERSION file is empty."
    exit 1
fi

log_info "Preparing Plachta ${PACKAGE_VERSION}"

STAGE_DIR="$(mktemp -d "${INSTALL_DIR}.stage.XXXXXX")"
install -m 0755 "$TMP_DIR/plachta" "$STAGE_DIR/plachta"
install -m 0644 "$TMP_DIR/VERSION" "$STAGE_DIR/VERSION"
cp -a "$TMP_DIR/modules" "$STAGE_DIR/modules"

if [[ ! -x "$STAGE_DIR/plachta" || ! -f "$STAGE_DIR/VERSION" || ! -d "$STAGE_DIR/modules" ]]; then
    log_error "Staged installation is incomplete."
    exit 1
fi

if [[ -e "$INSTALL_DIR" ]]; then
    BACKUP_DIR="${INSTALL_DIR}.backup.$$"
    mv "$INSTALL_DIR" "$BACKUP_DIR"
fi

if ! mv "$STAGE_DIR" "$INSTALL_DIR"; then
    log_error "Could not deploy the new installation."
    if [[ -n "$BACKUP_DIR" && -d "$BACKUP_DIR" && ! -e "$INSTALL_DIR" ]]; then
        mv "$BACKUP_DIR" "$INSTALL_DIR" || true
    fi
    exit 1
fi

STAGE_DIR=""
DEPLOYED=1

ln -sfn "$INSTALL_DIR/plachta" "$BIN_LINK"

log_info "Installation files deployed."
"$BIN_LINK" version
log_info "Plachta installation completed."
