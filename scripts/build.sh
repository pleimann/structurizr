#!/usr/bin/env bash

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")

env GOOS=darwin GOARCH=arm64 go build -ldflags "-s" -o strender-darwin-arm64 $SCRIPT_DIR/../

env GOOS=darwin GOARCH=arm64 go build -ldflags "-s" -o strender-darwin-amd64 $SCRIPT_DIR/../

env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o strender-linux-amd64 $SCRIPT_DIR/../

env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o strender-linux-arm64 $SCRIPT_DIR/../