#!/usr/bin/env bash

MODULES_DIR="${ROOT_DIR}/modules"

module_exists() {
    [[ -d "${MODULES_DIR}/$1" ]]
}

module_install() {
    local module="$1"

    if ! module_exists "$module"; then
        echo "Module '$module' not found."
        return 1
    fi

    bash "${MODULES_DIR}/${module}/install.sh"
}

module_verify() {
    local module="$1"

    if ! module_exists "$module"; then
        echo "Module '$module' not found."
        return 1
    fi

    bash "${MODULES_DIR}/${module}/verify.sh"
}

module_service() {
    local module="$1"

    if ! module_exists "$module"; then
        echo "Module '$module' not found."
        return 1
    fi

    bash "${MODULES_DIR}/${module}/service.sh"
}

module_config() {
    local module="$1"

    if ! module_exists "$module"; then
        echo "Module '$module' not found."
        return 1
    fi

    bash "${MODULES_DIR}/${module}/config.sh"
}

module_list() {
    for dir in "${MODULES_DIR}"/*; do
        [[ -d "$dir" ]] && basename "$dir"
    done
}