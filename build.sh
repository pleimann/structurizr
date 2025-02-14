env GOOS=darwin GOARCH=arm64 go build -ldflags "-s" -o strender-darwin-arm64 .

env GOOS=darwin GOARCH=arm64 go build -ldflags "-s" -o strender-darwin-amd64 .

env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o strender-linux-amd64 .

env GOOS=linux GOARCH=amd64 go build -ldflags "-s" -o strender-linux-arm64 .