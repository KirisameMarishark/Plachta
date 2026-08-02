#!/usr/bin/env bash

service_enable() {
    local service="$1"

    log_info "Enabling service: ${service}"

    systemctl enable "${service}"

    log_success "Service '${service}' enabled."
}

service_disable() {
    local service="$1"

    log_info "Disabling service: ${service}"

    systemctl disable "${service}"

    log_success "Service '${service}' disabled."
}

service_start() {
    local service="$1"

    log_info "Starting service: ${service}"

    systemctl start "${service}"

    log_success "Service '${service}' started."
}

service_stop() {
    local service="$1"

    log_info "Stopping service: ${service}"

    systemctl stop "${service}"

    log_success "Service '${service}' stopped."
}

service_restart() {
    local service="$1"

    log_info "Restarting service: ${service}"

    systemctl restart "${service}"

    log_success "Service '${service}' restarted."
}

service_status() {
    local service="$1"

    systemctl status "${service}"
}