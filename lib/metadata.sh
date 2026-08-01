#!/usr/bin/env bash

module_name() {
    basename "$1"
}

module_readme() {
    local module="$1"

    if [[ -f "${MODULES_DIR}/${module}/README.md" ]]; then
        cat "${MODULES_DIR}/${module}/README.md"
    else
        echo "No README found."
    fi
}

module_info() {
    local module="$1"

    if ! module_exists "$module"; then
        log_error "Module '$module' not found."
        return 1
    fi

    echo "Module : $module"
    echo "Path   : ${MODULES_DIR}/${module}"

    if [[ -f "${MODULES_DIR}/${module}/README.md" ]]; then
        echo "README : yes"
    else
        echo "README : no"
    fi

    if [[ -f "${MODULES_DIR}/${module}/install.sh" ]]; then
        echo "Install: yes"
    else
        echo "Install: no"
    fi

    if [[ -f "${MODULES_DIR}/${module}/verify.sh" ]]; then
        echo "Verify : yes"
    else
        echo "Verify : no"
    fi

    if [[ -f "${MODULES_DIR}/${module}/config.sh" ]]; then
        echo "Config : yes"
    else
        echo "Config : no"
    fi
}