//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestKubectlBuildkitBuildSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}
	createBuilder(t)

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
  kubectlBuildkit:
    build:
      file: dev/Dockerfile.dev
- image: simple-app-three
  path: assets/simple-app
  kubectlBuildkit:
    build:
      target: build-env
`)

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	out = regexp.MustCompile("rand-\\d+-\\d+").ReplaceAllString(out, "rand-REPLACED")

	expectedOut := `---
kind: Object
spec:
- image: kilnbase:rand-REPLACED-simple-app-two
- image: kilnbase:rand-REPLACED-simple-app-three
`

	if out != expectedOut {
		t.Fatalf("Expected >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestKubectlBuildkitBuildAndPushSuccessful(t *testing.T) {
	env := BuildEnv(t)

	if env.SkipWhenHTTPRegistry {
		fmt.Printf("This is a test that cannot run against HTTP registry; skipping.")
		return
	}

	createBuilder(t)

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
  kubectlBuildkit: {}
- image: docker.io/*username*/kilnbase-e2e-tests-build2
  path: assets/simple-app
  kubectlBuildkit: {}
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageDestinations
destinations:
- image: docker.io/*username*/kilnbase-e2e-tests-build
- image: docker.io/*username*/kilnbase-e2e-tests-build2
`)

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	// Expecting same sha256 multiple times
	out = strings.Replace(out, regexp.MustCompile("sha256:[a-z0-9]{64}").FindString(out), "SHA256-REPLACED", -1)

	expectedOut := env.WithRegistries(`---
kind: Object
spec:
- image: index.docker.io/*username*/kilnbase-e2e-tests-build@SHA256-REPLACED
- image: index.docker.io/*username*/kilnbase-e2e-tests-build2@SHA256-REPLACED
`)

	if out != expectedOut {
		t.Fatalf("Expected >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

// This function was created ensure that a builder was present in the cluster
// before we try to run these tests.
// If a builder does not exist the build command will try to create one.
// The problem is that we run multiple build commands in parallel which cause
// a race condition described in
