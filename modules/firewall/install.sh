#!/usr/bin/env bash
set -euo pipefail

require_root
require_debian
require_systemd

require_package nftables

ensure_service_enabled nftables
ensure_service_running nftables

log_success "Firewall installation completed."