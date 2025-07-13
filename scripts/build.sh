#!/usr/bin/env bash

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")
PROJECT_DIR=$SCRIPT_DIR/../

env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/strender-darwin-arm64 $PROJECT_DIR

env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/strender-darwin-amd64 $PROJECT_DIR

env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/strender-linux-amd64 $PROJECT_DIR

env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o bin/strender-linux-arm64 $PROJECT_DIR