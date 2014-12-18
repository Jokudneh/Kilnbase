#!/bin/bash

set -e -x -u

export GOTOOLCHAIN=auto

./hack/build.sh

export KILNBASE_BINARY_PATH="$PWD/kilnbase"

./hack/test.sh
