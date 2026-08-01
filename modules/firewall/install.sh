#!/usr/bin/env bash

set -euo pipefail

log_info "Checking operating system..."

if [[ ! -f /etc/os-release ]]; then
    log_error "Cannot determine operating system."
    exit 1
fi

source /etc/os-release

if [[ "${ID}" != "debian" ]]; then
    log_error "Only Debian is currently supported."
    exit 1
fi

log_success "Operating system: Debian"

log_info "Checking root privilege..."

if [[ "${EUID}" -ne 0 ]]; then
    log_error "Please run as root."
    exit 1
fi

log_success "Running as root."

log_info "Checking nftables..."

if command -v nft >/dev/null 2>&1; then
    log_success "nftables is already installed."
else
    log_warn "nftables not found."

    log_info "Installing nftables..."

    apt-get update

    apt-get install -y nftables

    log_success "nftables installed."
fi

log_info "Enabling nftables service..."

systemctl enable nftables

systemctl start nftables

log_success "Firewall service started."