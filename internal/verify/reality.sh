#!/usr/bin/env bash
set -euo pipefail

verify_reality() {
    local failed=0

    echo "Reality Verify"
    echo "------------------------------"

    # Config
    if [[ -f /etc/plachta/reality/config.json ]]; then
        echo "✔ Config exists"
    else
        echo "✘ Config missing"
        failed=1
    fi

    # Xray binary
    if command -v xray >/dev/null 2>&1; then
        echo "✔ Xray installed"
    else
        echo "✘ Xray not installed"
        failed=1
    fi

    # Service
    if systemctl is-active --quiet xray; then
        echo "✔ Xray running"
    else
        echo "✘ Xray stopped"
        failed=1
    fi

    # Port
    if ss -lnt | grep -q ':443 '; then
        echo "✔ Port 443 listening"
    else
        echo "✘ Port 443 not listening"
        failed=1
    fi

    # Config syntax
    if xray run -test -config /etc/plachta/reality/config.json >/dev/null 2>&1; then
        echo "✔ Config valid"
    else
        echo "✘ Config invalid"
        failed=1
    fi

    echo

    if [[ $failed -eq 0 ]]; then
        echo "Reality looks good."
        return 0
    fi

    echo "Reality has problems."
    return 1
}