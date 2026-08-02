#!/usr/bin/env bash
DEFAULT_CONFIG="${ROOT_DIR}/configs/default/plachta.conf"
CONFIG_DIR="${HOME}/.config/plachta"
CONFIG_FILE="${CONFIG_DIR}/config.conf"

config_init() {
    mkdir -p "${CONFIG_DIR}"

    if [[ ! -f "${CONFIG_FILE}" ]]; then
        cat > "${CONFIG_FILE}" <<EOF
# Plachta Configuration

DNS=1.1.1.1
MTU=1500
LOG_LEVEL=info
EOF
    fi
}

config_get() {
    local key="$1"
    grep "^${key}=" "${CONFIG_FILE}" | cut -d= -f2-
}

config_set() {
    local key="$1"
    local value="$2"

    if grep -q "^${key}=" "${CONFIG_FILE}"; then
        sed -i "s|^${key}=.*|${key}=${value}|" "${CONFIG_FILE}"
    else
        echo "${key}=${value}" >> "${CONFIG_FILE}"
    fi
}