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

    cat >/etc/plachta/reality/config.json <<EOF
    cat >/etc/plachta/reality/client.env <<EOF
    UUID=$uuid
    PUBLIC_KEY=$public_key
    SHORT_ID=$shortid
    SERVER_NAME=www.cloudflare.com
    PORT=443
    EOF
{
  "log": {
    "loglevel": "warning"
  },
  "inbounds": [
    {
      "listen": "0.0.0.0",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": [
          {
            "id": "$uuid",
            "flow": "xtls-rprx-vision"
          }
        ],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "tcp",
        "security": "reality",
        "realitySettings": {
          "show": false,
          "dest": "www.cloudflare.com:443",
          "xver": 0,
          "serverNames": [
            "www.cloudflare.com"
          ],
          "privateKey": "$private_key",
          "shortIds": [
            "$shortid"
          ]
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom"
    }
  ]
}
EOF

    cat >/etc/plachta/reality/client.env <<EOF
UUID=$uuid
PUBLIC_KEY=$public_key
SHORT_ID=$shortid
SERVER_NAME=www.cloudflare.com
PORT=443
EOF

    log_success "Reality config generated."
}
