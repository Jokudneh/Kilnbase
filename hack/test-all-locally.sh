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
  echo $env | grep -E 'export*' | awk '{print $2}' | sed 's/"//g'
done > $tempConfigFile

docker run \
--privileged \
--env-file $tempConfigFile \
-e DOCKER_API_VERSION=1.44 \
-e KILNBASE_E2E_SKIP_WHEN_HTTP_REGISTRY=${KILNBASE_E2E_SKIP_WHEN_HTTP_REGISTRY:-true} \
-e KILNBASE_E2E_DOCKERHUB_HOSTNAME=`minikube ip`:30777 \
-v ~/.config:/root/.config \
-v ~/.minikube:"$HOME/.minikube" \
-v ~/.kube:/root/.kube \
-v /etc/docker/:/host-etc-docker \
