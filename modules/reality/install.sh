#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/lib/common.sh"
source "${ROOT_DIR}/lib/logger.sh"
source "${ROOT_DIR}/lib/prereq.sh"
source "${ROOT_DIR}/lib/package.sh"

source "${ROOT_DIR}/internal/xray/download.sh"
source "${ROOT_DIR}/internal/xray/install.sh"
source "${ROOT_DIR}/internal/xray/service.sh"
source "${ROOT_DIR}/internal/reality/generate.sh"

MODULES_DIR="${ROOT_DIR}/modules"

source "${MODULES_DIR}/reality/config.sh"

require_root
require_debian
require_systemd

package_install jq

xray_download
xray_install
generate_reality_config
xray_install_service

log_success "Reality installation completed."