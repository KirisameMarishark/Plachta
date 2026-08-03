#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/lib/common.sh"
source "${ROOT_DIR}/lib/logger.sh"
source "${ROOT_DIR}/lib/prereq.sh"
source "${ROOT_DIR}/lib/package.sh"
source "${ROOT_DIR}/lib/service.sh"

MODULES_DIR="${ROOT_DIR}/modules"

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