#!/usr/bin/env bash

cmd_install() {

    if [[ -z "${2:-}" ]]; then
        echo "Usage:"
        echo "  plachta install <module>"
        exit 1
    fi

    log_info "Installing module: $2"

    module_install "$2"

    log_success "Install finished."
}