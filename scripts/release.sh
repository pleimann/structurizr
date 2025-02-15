./build.sh

gh release create $1 --draft --latest "strender-darwin-arm64#MacOS Arm" "strender-darwin-amd64#MacOS x64" "strender-linux-arm64#Linux Arm" "strender-linux-amd64#Linux x64"