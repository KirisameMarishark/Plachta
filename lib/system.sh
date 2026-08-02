#!/usr/bin/env bash

system_os() {
    source /etc/os-release
    echo "${ID}"
}

system_version() {
    source /etc/os-release
    echo "${VERSION_ID}"
}

system_arch() {
    uname -m
}

system_kernel() {
    uname -r
}

system_hostname() {
    hostname
}

system_ipv4() {
    ip -4 route get 1.1.1.1 2>/dev/null | awk '{print $7; exit}'
}

system_ipv6() {
    ip -6 route get 2606:4700:4700::1111 2>/dev/null | awk '{print $7; exit}'
}

system_memory() {
    free -m | awk '/Mem:/ {print $2}'
}

system_cpu() {
    nproc
}