#!/usr/bin/env bash

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")
PROJECT_DIR=$SCRIPT_DIR/../

$SCRIPT_DIR/build.sh

ARCH=$(arch)
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')

cp $PROJECT_DIR/bin/strender-$PLATFORM-$ARCH ~/bin/strender