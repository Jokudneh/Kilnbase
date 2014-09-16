#!/bin/bash

set -e -x -u

git --version || (echo "Missing git binary (used by tests)" && exit 1)
