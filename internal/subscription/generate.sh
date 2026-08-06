#!/usr/bin/env bash
set -euo pipefail

source "${ROOT_DIR}/internal/reality/uri.sh"

generate_subscription() {

    mkdir -p /etc/plachta/subscription

    reality_generate_uri >/etc/plachta/subscription/reality.txt

}