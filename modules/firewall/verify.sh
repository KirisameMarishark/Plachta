#!/usr/bin/env bash

set -euo pipefail

log_info "Checking nftables..."

if ! command -v nft >/dev/null 2>&1; then
    log_error "nft command not found."
    exit 1
fi

log_success "nft command found."

log_info "Checking nftables service..."

if systemctl is-enabled nftables >/dev/null 2>&1; then
    log_success "nftables service enabled."
else
    log_warn "nftables service is not enabled."
fi

if systemctl is-active nftables >/dev/null 2>&1; then
    log_success "nftables service running."
else
    log_warn "nftables service is not running."
fi

log_info "Listing ruleset..."

nft list ruleset