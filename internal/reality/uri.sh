#!/usr/bin/env bash
set -euo pipefail

source "${ROOT_DIR}/internal/reality/read.sh"

reality_uri() {

    local server_ip
    server_ip="$(curl -4 -fsSL https://api.ipify.org)"

    printf \
'vless://%s@%s:%s?type=tcp&security=reality&pbk=%s&sid=%s&sni=%s&fp=chrome&flow=xtls-rprx-vision&encryption=none#Plachta-Reality' \
        "$(reality_uuid)" \
        "$server_ip" \
        "$(reality_port)" \
        "$(reality_public_key)" \
        "$(reality_short_id)" \
        "$(reality_server_name)"
}