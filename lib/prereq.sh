#!/usr/bin/env bash

require_root() {
    if [[ "${EUID}" -ne 0 ]]; then
        log_error "Please run as root."
        exit 1
    fi

    log_success "Root privilege confirmed."
}

require_debian() {
    if [[ ! -f /etc/os-release ]]; then
        log_error "Cannot determine operating system."
        exit 1
    fi

    source /etc/os-release

    if [[ "${ID}" != "debian" ]]; then
        log_error "Only Debian is supported."
        exit 1
    fi

    log_success "Operating system: Debian"
}

require_systemd() {
    if ! command -v systemctl >/dev/null 2>&1; then
        log_error "systemd is required."
        exit 1
    fi

    log_success "systemd detected."
}

require_command() {
    local cmd="$1"

    if ! command -v "$cmd" >/dev/null 2>&1; then
        log_error "Command '$cmd' not found."
        return 1
    fi

    log_success "Command '$cmd' found."
}

require_package() {
    local pkg="$1"

    if dpkg -s "$pkg" >/dev/null 2>&1; then
        log_success "Package '$pkg' already installed."
        return 0
    fi

    log_warn "Installing package '$pkg'..."

    apt-get update

    apt-get install -y "$pkg"

    log_success "Package '$pkg' installed."
}

ensure_service_enabled() {
    service_enable "$1"
}

ensure_service_running() {
    service_start "$1"
}