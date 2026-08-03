#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/lib/common.sh"
source "${ROOT_DIR}/lib/logger.sh"
source "${ROOT_DIR}/lib/prereq.sh"

source "${ROOT_DIR}/internal/xray/download.sh"
source "${ROOT_DIR}/internal/xray/install.sh"

require_root
require_debian
require_systemd

xray_download
xray_install

log_success "Reality installation completed."