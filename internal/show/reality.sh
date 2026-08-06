#!/usr/bin/env bash
set -euo pipefail

source "${ROOT_DIR}/internal/reality/read.sh"
source "${ROOT_DIR}/internal/reality/uri.sh"

show_reality() {

    if [[ ! -f /etc/plachta/reality/config.json ]]; then
        echo "Reality is not installed."
        exit 1
    fi

    local server_ip
    server_ip="$(curl -4 -fsSL https://api.ipify.org)"

    local uri
    uri="$(reality_uri)"

    echo
    echo "Reality"
    echo "------------------------------"
    echo "Server     : ${server_ip}"
    echo "Port       : $(reality_port)"
    echo "UUID       : $(reality_uuid)"
    echo "Public Key : $(reality_public_key)"
    echo "ServerName : $(reality_server_name)"
    echo "Short ID   : $(reality_short_id)"
    echo

    echo "Import URI:"
    echo
    echo "$uri"
    echo

    echo "QR Code"
    echo "------------------------------"

    if command -v qrencode >/dev/null 2>&1; then

        qrencode -t ANSIUTF8 "$uri"

        qrencode \
            -o /etc/plachta/reality/reality.png \
            -s 8 \
            -m 2 \
            "$uri"

        echo
        echo "PNG QR Code"
        echo "------------------------------"
        echo "/etc/plachta/reality/reality.png"

    else
        echo "qrencode not installed."
    fi
}