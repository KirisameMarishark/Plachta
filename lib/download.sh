#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

source "${ROOT_DIR}/lib/logger.sh"
download_has_curl() {
    command -v curl >/dev/null 2>&1
}

download_has_wget() {
    command -v wget >/dev/null 2>&1
}

download_fetch() {
    local url="$1"
    local output="$2"

    if download_has_curl; then
        log_info "Downloading with curl..."
        curl -L --fail --retry 3 --retry-delay 2 -o "$output" "$url"
        return $?
    fi

    if download_has_wget; then
        log_info "Downloading with wget..."
        wget -O "$output" "$url"
        return $?
    fi

    die "Neither curl nor wget is available."
}