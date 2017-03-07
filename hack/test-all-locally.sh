#!/bin/bash
set -e -x -u

go clean -testcache
./hack/build.sh
export KILNBASE_BINARY_PATH="${KILNBASE_BINARY_PATH:-$PWD/kilnbase}"

function build_test_deps() {
cat <<EOF | $KILNBASE_BINARY_PATH -f - | grep 'image:' | awk '{print $2}'
image: test-dependencies
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
sources:
- image: test-dependencies
  path: .
  docker:
    build:
      pull: true
      noCache: false
      file: hack/Dockerfile.dev
EOF
}

image_name=$(build_test_deps)

tempConfigFile=$(mktemp)
trap "rm -f $tempConfigFile" EXIT

minikube docker-env | while read env; do
