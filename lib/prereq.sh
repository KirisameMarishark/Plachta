#!/usr/bin/env bash

require_root() {
    if [[ "${EUID}" -ne 0 ]]; then
        die "Please run as root."
        exit 1
    fi

    log_success "Root privilege confirmed."
}

require_debian() {
    if [[ ! -f /etc/os-release ]]; then
        die "Cannot determine operating system."
    fi

    source /etc/os-release

    if [[ "${ID}" != "debian" ]]; then
        die "Only Debian is supported."   
    fi

    log_success "Operating system: Debian"
}

require_systemd() {
    if ! command -v systemctl >/dev/null 2>&1; then
        die "systemd is required."    
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
    package_install "$1"
}
ensure_service_enabled() {
    service_enable "$1"
}

ensure_service_running() {
    service_start "$1"
}