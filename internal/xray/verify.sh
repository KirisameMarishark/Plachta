#!/usr/bin/env bash

xray_verify() {

    if command -v xray >/dev/null 2>&1; then
        xray version
        return 0
    fi

    die "Xray not installed."

}