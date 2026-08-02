#!/usr/bin/env bash

XRAY_VERSION=""

xray_version_set() {
    XRAY_VERSION="$1"
}

xray_version_get() {
    echo "${XRAY_VERSION}"
}