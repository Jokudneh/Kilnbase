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

