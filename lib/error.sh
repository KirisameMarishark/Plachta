#!/usr/bin/env bash

die() {
    local message="$1"

    log_error "$message"

    exit 1
}

warn() {
    local message="$1"

    log_warn "$message"
}

success() {
    local message="$1"

    log_success "$message"
}

info() {
    local message="$1"

    log_info "$message"
}