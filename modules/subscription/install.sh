#!/usr/bin/env bash
set -euo pipefail

source "${ROOT_DIR}/internal/subscription/generate.sh"

generate_subscription

log_success "Subscription generated."