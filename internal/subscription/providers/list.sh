#!/usr/bin/env bash
set -euo pipefail

subscription_provider_list() {

    local provider_dir

    provider_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

    find "$provider_dir" \
        -maxdepth 1 \
        -type f \
        -name "*.sh" \
        ! -name "list.sh" \
        -printf "%f\n" \
    | sed 's/\.sh$//' \
    | sort

}