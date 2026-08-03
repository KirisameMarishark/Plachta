#!/usr/bin/env bash

github_latest_release() {

    local repo="$1"

    curl -fsSL \
        -H "Accept: application/vnd.github+json" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "https://api.github.com/repos/${repo}/releases/latest"

}