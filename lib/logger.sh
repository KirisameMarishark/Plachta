#!/usr/bin/env bash

log_time() {
    date "+%Y-%m-%d %H:%M:%S"
}

log_info() {
    echo "[INFO ] $(log_time) $*"
}

log_warn() {
    echo "[WARN ] $(log_time) $*"
}

log_error() {
    echo "[ERROR] $(log_time) $*" >&2
}

log_success() {
    echo "[ OK  ] $(log_time) $*"
}