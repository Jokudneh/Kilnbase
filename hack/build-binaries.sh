#!/bin/bash

set -e -x -u

function get_latest_git_tag {
  git describe --tags | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+'
}

VERSION="${1:-`get_latest_git_tag`}"

# makes builds reproducible
export CGO_ENABLED=0
LDFLAGS="-X github.com/Jokudneh/Kilnbase/pkg/kilnbase/version.Version=$VERSION"

GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kilnbase-darwin-amd64 ./cmd/kilnbase/...
GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kilnbase-darwin-arm64 ./cmd/kilnbase/...
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kilnbase-linux-amd64 ./cmd/kilnbase/...
GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -trimpath -o kilnbase-linux-arm64 ./cmd/kilnbase/...
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -trimpath -o kilnbase-windows-amd64.exe ./cmd/kilnbase/...
