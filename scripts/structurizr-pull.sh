#!/usr/bin/env bash

SCRIPT_DIR=$(dirname -- "$(readlink -f "${BASH_SOURCE}")")

FILES="js/structurizr.js
js/structurizr-content.js
js/structurizr-diagram.js
js/structurizr-ui.js
js/structurizr-util.js
js/structurizr-workspace.js
css/structurizr.css
css/structurizr-diagram.css"

for file in $FILES; do
  echo "Downloading $file..."
  curl -o "$SCRIPT_DIR/../frontend/$file" "https://raw.githubusercontent.com/structurizr/ui/refs/heads/main/src/$file"
done
