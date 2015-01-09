//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRelocateSuccessful(t *testing.T) {
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

	out, _ := kilnbase.RunWithOpts([]string{"relocate", "-f", "-", "--repository", env.WithRegistries("docker.io/*username*/kilnbase-test-relocate")}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := env.WithRegistries(`---
kind: Object
spec:
- image: index.docker.io/*username*/kilnbase-test-relocate@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
- image: index.docker.io/*username*/kilnbase-test-relocate@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
`)

	if out != expectedOut {
		t.Fatalf("Expected unpackage output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestRelocateSuccessfulWithManyImages(t *testing.T) {
	env := BuildEnv(t)

	if env.SkipStressTests {
