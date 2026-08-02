#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/lib/github.sh"

XRAY_REPO="XTLS/Xray-core"

xray_latest_version() {

    github_latest_release "${XRAY_REPO}" \
        | grep '"tag_name":' \
        | head -n1 \
        | cut -d '"' -f4

}