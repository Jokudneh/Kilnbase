#!/bin/bash

set -e -x -u

LATEST_GIT_TAG=$(git describe --tags | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
VERSION="${1:-$LATEST_GIT_TAG}"

# makes builds reproducible
export CGO_ENABLED=0
LDFLAGS="-X github.com/Jokudneh/Kilnbase/pkg/kilnbase/version.Version=$VERSION"
