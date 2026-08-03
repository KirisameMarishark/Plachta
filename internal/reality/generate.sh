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

reality_private_key() {
    xray x25519 | awk -F': ' '/^PrivateKey:/ {print $2}'
}

reality_public_key() {
    xray x25519 | awk -F': ' '/^Password \(PublicKey\):/ {print $2}'
}

reality_hash() {
    xray x25519 | awk -F': ' '/^Hash32:/ {print $2}'
}