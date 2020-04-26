//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"regexp"
	"strings"
	"testing"
)

func TestDockerBuildxBuildSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := env.WithRegistries(`
kind: Object
spec:
- image: simple-app-two
- image: simple-app-three
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Sources
sources:
- image: simple-app-two
  path: assets/simple-app
  docker:
    buildx:
      file: dev/Dockerfile.dev
- image: simple-app-three
  path: assets/simple-app
  docker:
    buildx:
      target: build-env
`)

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	out = regexp.MustCompile("sha256\\-[a-z0-9]{64}").ReplaceAllString(out, "SHA256-REPLACED")

	expectedOut := `---
kind: Object
spec:
- image: kilnbase:simple-app-two-SHA256-REPLACED
- image: kilnbase:simple-app-three-SHA256-REPLACED
`

	if out != expectedOut {
		t.Fatalf("Expected >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestDockerBuildxBuildAndPushSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := env.WithRegistries(`
kind: Object
spec:
- image: docker.io/*username*/kilnbase-e2e-tests-build
- image: docker.io/*username*/kilnbase-e2e-tests-build2
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Sources
sources:
- image: docker.io/*username*/kilnbase-e2e-tests-build
  path: assets/simple-app
  docker:
    buildx: {}
- image: docker.io/*username*/kilnbase-e2e-tests-build2
  path: assets/simple-app
  docker:
    buildx:
      # try out multi platform build
      rawOptions: ["--platform=linux/amd64,linux/arm64,linux/arm/v7"]
---
apiVersion: kilnbase.k14s.io/v1alpha1
