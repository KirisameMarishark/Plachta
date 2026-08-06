#!/usr/bin/env bash
set -euo pipefail

source "${ROOT_DIR}/internal/reality/read.sh"

validate_reality() {

    load_reality_config

    [[ -f /etc/plachta/reality/config.json ]] || return 1

    command -v jq >/dev/null 2>&1 || return 1

    command -v xray >/dev/null 2>&1 || return 1

    jq empty /etc/plachta/reality/config.json >/dev/null 2>&1 || return 1

    [[ -n "$REALITY_UUID" ]] || return 1

    [[ -n "$REALITY_PRIVATE_KEY" ]] || return 1

    [[ -n "$REALITY_SHORT_ID" ]] || return 1

    [[ "$REALITY_SECURITY" == "reality" ]] || return 1

    systemctl is-active --quiet xray || return 1

    ss -lntp | grep -q xray || return 1

    return 0
}