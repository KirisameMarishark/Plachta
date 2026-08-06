#!/usr/bin/env bash

cmd_verify() {

    if [[ -z "${2:-}" ]]; then
        echo "Usage:"
        echo "  plachta verify reality"
        exit 1
    fi

    case "$2" in

        reality)

            verify_reality
            ;;

        *)

            module_verify "$2"
            ;;
    esac

}