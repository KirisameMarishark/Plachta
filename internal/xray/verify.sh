#!/usr/bin/env bash

xray_verify() {

    if command -v xray >/dev/null 2>&1; then
        success "Xray detected."
        return 0
    fi

    die "Xray is not installed."

}