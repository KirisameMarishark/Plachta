#!/usr/bin/env bash
set -euo pipefail

generate_reality_config() {

    mkdir -p /etc/plachta/reality

    uuid="$(reality_generate_uuid)"

    keypair="$(reality_generate_keypair)"

    private_key="$(printf '%s\n' "$keypair" | awk -F': ' '/^PrivateKey:/ {print $2}')"

    public_key="$(printf '%s\n' "$keypair" | awk -F': ' '/^Password \(PublicKey\):/ {print $2}')"

    hash32="$(printf '%s\n' "$keypair" | awk -F': ' '/^Hash32:/ {print $2}')"

    shortid="$(reality_generate_shortid)"

    echo "UUID=$uuid"
    echo "PRIVATE=$private_key"
    echo "PUBLIC=$public_key"
    echo "HASH=$hash32"
    echo "SHORTID=$shortid"

    cat >/etc/plachta/reality/config.json <<EOF
{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [],
  "outbounds": [
    {
      "protocol": "freedom"
    }
  ]
}
EOF

    log_success "Reality config generated."
}