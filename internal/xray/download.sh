#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/lib/download.sh"
source "${ROOT_DIR}/internal/xray/version.sh"

xray_download() {

    local version
    version="$(xray_latest_version)"

    local file="Xray-linux-64.zip"

    local url="https://github.com/XTLS/Xray-core/releases/download/${version}/${file}"

    mkdir -p "${ROOT_DIR}/runtime"

    download_fetch "$url" "${ROOT_DIR}/runtime/${file}"

}