#!/usr/bin/env bash
set -euo pipefail

verify_reality() {
    local pass=0
    local fail=0

    CONFIG="/etc/plachta/reality/config.json"

    echo "Reality Verify"
    echo "------------------------------"

    ok() {
        echo "✔ $1"
        ((pass++))
    }

    bad() {
        echo "✘ $1"
        ((fail++))
    }

    # jq
    if command -v jq >/dev/null 2>&1; then
        ok "jq installed"
    else
        bad "jq missing"
    fi

    # Config
    if [[ -f "$CONFIG" ]]; then
        ok "Config exists"
    else
        bad "Config missing"
    fi

    # JSON
    if jq empty "$CONFIG" >/dev/null 2>&1; then
        ok "JSON valid"
    else
        bad "JSON invalid"
    fi

    # UUID
    if [[ -n "$(jq -r '.inbounds[0].settings.clients[0].id' "$CONFIG")" ]]; then
        ok "UUID found"
    else
        bad "UUID missing"
    fi

    # PrivateKey
    if [[ -n "$(jq -r '.inbounds[0].streamSettings.realitySettings.privateKey' "$CONFIG")" ]]; then
        ok "PrivateKey found"
    else
        bad "PrivateKey missing"
    fi

    # ShortID
    if [[ -n "$(jq -r '.inbounds[0].streamSettings.realitySettings.shortIds[0]' "$CONFIG")" ]]; then
        ok "ShortID found"
    else
        bad "ShortID missing"
    fi

    # PrivateKey 可用性
    if xray x25519 -i "$(jq -r '.inbounds[0].streamSettings.realitySettings.privateKey' "$CONFIG")" >/dev/null 2>&1; then
        ok "PrivateKey usable"
    else
        bad "PrivateKey invalid"
    fi

    # Reality
    if [[ "$(jq -r '.inbounds[0].streamSettings.security' "$CONFIG")" == "reality" ]]; then
        ok "Reality inbound"
    else
        bad "Reality inbound missing"
    fi

    # Xray
    if command -v xray >/dev/null 2>&1; then
        ok "Xray installed"
    else
        bad "Xray missing"
    fi

    # Service
    if systemctl is-active --quiet xray; then
        ok "Xray running"
    else
        bad "Xray stopped"
    fi

    # Port
    if ss -lntp | grep -q 'xray'; then
        ok "Port 443 listening"
    else
        bad "Port 443 not listening"
    fi

    echo
    echo "Result"
    echo "------------------------------"
    echo "Passed : $pass"
    echo "Failed : $fail"

    if [[ -f /etc/plachta/reality/client.env ]]; then
        source /etc/plachta/reality/client.env

        server_ip="$(curl -4 -fsSL https://api.ipify.org 2>/dev/null || echo "Unknown")"

        uri="vless://${UUID}@${server_ip}:${PORT}?type=tcp&security=reality&pbk=${PUBLIC_KEY}&sid=${SHORT_ID}&sni=${SERVER_NAME}&fp=chrome&flow=xtls-rprx-vision&encryption=none#Plachta-Reality"

        echo
        echo "Current Server"
        echo "------------------------------"
        echo "IP   : $server_ip"
        echo "Port : $PORT"

        echo
        echo "Import URI"
        echo "------------------------------"
        echo "$uri"
    fi

    if [[ $fail -eq 0 ]]; then
        echo
        echo "Reality looks good."
        return 0
    fi

    echo
    echo "Reality has problems."
    return 1
}