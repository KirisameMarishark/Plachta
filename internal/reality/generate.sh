#!/usr/bin/env bash
set -euo pipefail

reality_generate_uuid() {
    xray uuid
}

reality_generate_keypair() {
    xray x25519
}

reality_generate_shortid() {
    openssl rand -hex 8
}