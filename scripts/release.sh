#!/usr/bin/env bash

if [[ -z "$1" ]]; then
    echo "ERROR: Must provide release number argument" 1>&2
    exit 1
fi

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")
PROJECT_DIR=$SCRIPT_DIR/..

$SCRIPT_DIR/build.sh

gh release create $1 --draft --latest "$PROJECT_DIR/bin/strender-darwin-arm64#MacOS Arm64" "$PROJECT_DIR/bin/strender-darwin-amd64#MacOS x64" "$PROJECT_DIR/bin/strender-linux-arm64#Linux Arm64" "$PROJECT_DIR/bin/strender-linux-amd64#Linux x64"