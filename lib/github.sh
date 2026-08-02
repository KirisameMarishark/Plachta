#!/usr/bin/env bash

github_latest_release() {

    local repo="$1"

    curl -fsSL "https://api.github.com/repos/${repo}/releases/latest"

}