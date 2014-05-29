#!/bin/bash

set -e -x -u

export GOTOOLCHAIN=auto

./hack/build.sh

