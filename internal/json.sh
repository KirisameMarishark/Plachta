#!/usr/bin/env bash
set -euo pipefail

json_get() {
    local file="$1"
    local query="$2"

    jq -r "$query" "$file"
}