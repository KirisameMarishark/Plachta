#!/usr/bin/env bash
set -euo pipefail

subscription_add_provider() {

    local provider="$1"

    if command -v "$provider" >/dev/null 2>&1; then
        "$provider"
    fi

}