#!/usr/bin/env bash

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")
PROJECT_DIR=$SCRIPT_DIR/../

env GOOS=darwin GOARCH=arm64 go build -ldflags "-s" -o bin/strender-darwin-arm64 $PROJECT_DIR

env GOOS=darwin GOARCH=arm64 go build -ldflags "-s" -o bin/strender-darwin-amd64 $PROJECT_DIR

env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o bin/strender-linux-amd64 $PROJECT_DIR

env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o bin/strender-linux-arm64 $PROJECT_DIR