#!/usr/bin/env bash
set -euo pipefail

reality_uri() {

    local env_file="/etc/plachta/reality/client.env"

    [[ -f "$env_file" ]] || return 1

    source "$env_file"

    local server_ip
    server_ip="$(curl -4 -fsSL https://api.ipify.org 2>/dev/null)"

    printf "vless://%s@%s:%s?type=tcp&security=reality&pbk=%s&sid=%s&sni=%s&fp=chrome&flow=xtls-rprx-vision&encryption=none#Plachta-Reality" \
        "$UUID" \
        "$server_ip" \
        "$PORT" \
        "$PUBLIC_KEY" \
        "$SHORT_ID" \
        "$SERVER_NAME"
}