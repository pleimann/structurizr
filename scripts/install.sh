#!/usr/bin/env bash

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")

$SCRIPT_DIR/build.sh

cp $SCRIPT_DIR/../strender-darwin-arm64 ~/bin/strender