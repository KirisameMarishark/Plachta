#!/usr/bin/env bash
set -euo pipefail

REALITY_CONFIG="/etc/plachta/reality/config.json"

reality_uuid() {
    jq -r '.inbounds[0].settings.clients[0].id' "$REALITY_CONFIG"
}

reality_private_key() {
    jq -r '.inbounds[0].streamSettings.realitySettings.privateKey' "$REALITY_CONFIG"
}

reality_short_id() {
    jq -r '.inbounds[0].streamSettings.realitySettings.shortIds[0]' "$REALITY_CONFIG"
}

reality_server_name() {
    jq -r '.inbounds[0].streamSettings.realitySettings.serverNames[0]' "$REALITY_CONFIG"
}

reality_port() {
    jq -r '.inbounds[0].port' "$REALITY_CONFIG"
}

reality_public_key() {
    xray x25519 -i "$(reality_private_key)" \
        | awk -F': ' '/PublicKey/ {print $2}'
}