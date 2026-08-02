#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

source "${ROOT_DIR}/internal/xray/download.sh"

xray_install() {

    xray_download

    mkdir -p "${ROOT_DIR}/runtime/xray"

    unzip -o \
        "${ROOT_DIR}/runtime/Xray-linux-64.zip" \
        -d "${ROOT_DIR}/runtime/xray"

}