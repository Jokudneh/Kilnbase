#!/bin/bash

set -e -x -u

export KILNBASE_E2E_DOCKERHUB_USERNAME=minikube-tests
export KILNBASE_E2E_DOCKERHUB_HOSTNAME=${KILNBASE_E2E_DOCKERHUB_HOSTNAME:-`minikube ip`}
export DOCKER_API_VERSION=1.44
# uncomment to disable stress tests
# export KILNBASE_E2E_SKIP_STRESS_TESTS=true
