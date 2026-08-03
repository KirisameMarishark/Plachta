#!/usr/bin/env bash
set -euo pipefail

generate_reality_config() {
    mkdir -p /etc/plachta/reality

    cat >/etc/plachta/reality/config.json <<'EOF'
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