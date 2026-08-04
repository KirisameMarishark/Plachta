#!/usr/bin/env bash
set -euo pipefail

show_reality() {

    local env_file="/etc/plachta/reality/client.env"

    if [[ ! -f "$env_file" ]]; then
        echo "Reality is not installed."
        exit 1
    fi

    source "$env_file"

    local server_ip
    server_ip="$(curl -4 -fsSL https://api.ipify.org)"

    echo
    echo "Reality"
    echo "------------------------------"
    echo "Server     : ${server_ip}"
    echo "Port       : ${PORT}"
    echo "UUID       : ${UUID}"
    echo "Public Key : ${PUBLIC_KEY}"
    echo "ServerName : ${SERVER_NAME}"
    echo "Short ID   : ${SHORT_ID}"
    echo

    echo "Import URI:"
    echo
    echo "vless://${UUID}@${server_ip}:${PORT}?type=tcp&security=reality&pbk=${PUBLIC_KEY}&sid=${SHORT_ID}&sni=${SERVER_NAME}&fp=chrome&flow=xtls-rprx-vision&encryption=none#Plachta-Reality"
}