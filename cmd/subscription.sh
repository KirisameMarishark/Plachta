#!/usr/bin/env bash
set -euo pipefail

cmd_subscription() {

    case "${2:-generate}" in

        generate)

            generate_subscription
            ;;

        *)

            echo "Usage:"
            echo "  plachta subscription generate"
            return 1
            ;;

    esac

}