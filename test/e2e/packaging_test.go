//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestPkgUnpkgSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	// redis:5.0.4
	input := `
kind: Object
spec:
# references image index
- image: index.docker.io/library/redis@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
# references plain image
- image: gcr.io/cloud-builders/gcs-fetcher@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
`

	path := "/tmp/kilnbase-test-pkg-unpkg-successful"
	defer os.RemoveAll(path)

	out, _ := kilnbase.RunWithOpts([]string{"package", "-f", "-", "--output", path}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := ""

	if out != expectedOut {
		t.Fatalf("Expected package output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}

	out, _ = kilnbase.RunWithOpts([]string{"unpackage", "-f", "-", "--input", path, "--repository", env.WithRegistries("docker.io/*username*/kilnbase-test-pkg-unpkg")}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut = env.WithRegistries(`---
kind: Object
spec:
- image: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
- image: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
`)

	if out != expectedOut {
		t.Fatalf("Expected unpackage output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestPkgUnpkgLockSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
overrides:
# ignored because it's not preresolved
- image: gcs-fetcher
  newImage: gcr.io/cloud-builders/gcs-fetcher@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
- image: redis
  newImage: index.docker.io/library/redis@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
  preresolved: true
`

	path := "/tmp/kilnbase-test-pkg-unpkg-with-lock-successful"
	defer os.RemoveAll(path)

	relocatedLockPath := path + "-relocated"
	defer os.RemoveAll(relocatedLockPath)

	out, _ := kilnbase.RunWithOpts([]string{"package", "-f", "-", "--output", path}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := ""

	if out != expectedOut {
		t.Fatalf("Expected package output >>>%s<<< to match >>>%s<<<", out, expectedOut)
