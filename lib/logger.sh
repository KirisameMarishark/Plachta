#!/usr/bin/env bash
#
# Plachta Logger Library
#

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

########################################
# Logger Configuration
########################################

LOG_LEVEL="${LOG_LEVEL:-INFO}"

########################################
# Internal
########################################

_log() {

    local level="$1"
    shift

    local message="$*"

    local timestamp
    timestamp="$(date '+%Y-%m-%d %H:%M:%S')"

    echo "[$timestamp] [$level] $message"
}

########################################
# Public APIs
########################################

logger_info() {

    _log INFO "$@"

}

logger_warn() {

    _log WARN "$@"

}

logger_error() {

    _log ERROR "$@"

}

logger_success() {

    _log SUCCESS "$@"

}