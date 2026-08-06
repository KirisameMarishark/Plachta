#!/usr/bin/env bash
set -euo pipefail

generate_subscription() {

    local output_dir="/etc/plachta/subscription"
    local output_file="${output_dir}/sub.txt"

    mkdir -p "$output_dir"

    : > "$output_file"

    while read -r provider; do

        subscription_add_provider "${provider}_uri" >> "$output_file"

    done < <(subscription_provider_list)

    log_success "Subscription generated:"
    echo "$output_file"

}