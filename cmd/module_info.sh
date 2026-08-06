#!/usr/bin/env bash

cmd_module_info() {

    if [[ -z "${2:-}" ]]; then
        echo "Usage:"
        echo "  plachta module-info <module>"
        exit 1
    fi

    module_info "$2"

}