#!/usr/bin/env bash
#
# Plachta - Common Library
# Version : 0.1.0
# Author  : Plachta Project
#

set -o errexit
set -o nounset
set -o pipefail

########################################
# Project Information
########################################

readonly PLACHTA_NAME="Plachta"
readonly PLACHTA_VERSION="0.1.0"

########################################
# Colors
########################################

readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

########################################
# Logging
########################################

log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

log_success() {
    echo -e "${GREEN}[ OK ]${NC} $*"
}

########################################
# Environment Checks
########################################

require_root() {
    if [[ "${EUID}" -ne 0 ]]; then
        log_error "Please run this command as root."
        exit 1
    fi
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

require_command() {
    if ! command_exists "$1"; then
        log_error "Required command not found: $1"
        exit 1
    fi
}

########################################
# File Helpers
########################################

backup_file() {

    local file="$1"

    if [[ -f "$file" ]]; then
        cp "$file" "${file}.bak"
        log_info "Backup created: ${file}.bak"
    fi
}

restore_file() {

    local file="$1"

    if [[ -f "${file}.bak" ]]; then
        mv "${file}.bak" "$file"
        log_success "Backup restored."
    fi
}