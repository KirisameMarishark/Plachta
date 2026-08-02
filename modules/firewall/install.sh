#!/usr/bin/env bash
set -euo pipefail

require_root
require_debian
require_systemd

require_package nftables

source "${MODULES_DIR}/firewall/config.sh"

generate_firewall_config

ensure_service_enabled nftables
ensure_service_running nftables
nft -f /etc/nftables.conf

log_success "Firewall rules loaded."

log_success "Firewall installation completed."