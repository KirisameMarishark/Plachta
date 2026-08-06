#!/usr/bin/env bash

cmd_config() {

    case "${2:-show}" in

        show)
            cat "${CONFIG_FILE}"
            ;;

        get)
            config_get "$3"
            ;;

        set)
            config_set "$3" "$4"
            log_success "Configuration updated."
            ;;

        *)
            echo "Usage:"
            echo "  plachta config show"
            echo "  plachta config get <key>"
            echo "  plachta config set <key> <value>"
            ;;
    esac

}