#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/lib/logger.sh"
source "${ROOT_DIR}/lib/service.sh"

xray_install_service() {

    cat >/etc/systemd/system/xray.service <<'EOF'
[Unit]
Description=Xray Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/xray run -config /etc/plachta/reality/config.json
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload

    ensure_service_enabled xray
    ensure_service_running xray
}