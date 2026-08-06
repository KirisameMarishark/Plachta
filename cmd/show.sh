#!/usr/bin/env bash

cmd_show() {

    case "${2:-}" in

        reality)

            show_reality
            ;;

        *)

            echo "Usage:"
            echo "  plachta show reality"
            ;;
    esac

}