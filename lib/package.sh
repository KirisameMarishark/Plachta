#!/usr/bin/env bash

package_update() {
    log_info "Updating package index..."

    apt-get update

    log_success "Package index updated."
}

package_install() {
    local package="$1"

    if dpkg -s "$package" >/dev/null 2>&1; then
        log_success "Package '$package' already installed."
        return 0
    fi

    package_update

    log_info "Installing package: ${package}"

    apt-get install -y "${package}"

    log_success "Package '${package}' installed."
}

package_remove() {
    local package="$1"

    log_info "Removing package: ${package}"

    apt-get remove -y "${package}"

    log_success "Package '${package}' removed."
}