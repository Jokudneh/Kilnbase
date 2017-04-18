//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"regexp"
	"strings"
	"testing"
)

func TestDockerBuildSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := env.WithRegistries(`
kind: Object
spec:
- image: docker.io/*username*/kilnbase-e2e-tests-build
- image: simple-app-two
- image: simple-app-two-overriden
- image: simple-app-three
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Sources
sources:
- image: docker.io/*username*/kilnbase-e2e-tests-build
  path: assets/simple-app
- image: simple-app-two
  path: assets/simple-app
  docker:
    build:
      file: dev/Dockerfile.dev
- image: simple-app-three
  path: assets/simple-app
  docker:
    build:
      target: build-env
# 'unused' should not be built
- image: unused
  path: invalid-dir
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageOverrides
overrides:
- image: simple-app-two-overriden
  newImage: simple-app-two
`)

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	// kilnbase:192-168-99-100-30777-minikube-tests-kilnbase-e2e-tests-SHA256-REPLACED
	out = regexp.MustCompile("sha256\\-[a-z0-9]{64}").ReplaceAllString(out, "SHA256-REPLACED")
	out = regexp.MustCompile("kilnbase:(.+)-kilnbase-e2e(\\-.*)-SHA256-REPLACED").ReplaceAllString(out, "kilnbase:img-title-SHA256-REPLACED")

	expectedOut := `---
