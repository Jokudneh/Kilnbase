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

func TestBazelBuildAndPushSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	assetPath := "assets/simple-app"
	secondAssetPath := "assets/simple-app-2"
	if err := exec.Command("cp", "-r",
		assetPath, secondAssetPath).Run(); err != nil {
		t.Fatalf("failed to copy %s to %s: %v", assetPath, secondAssetPath, err)
	}
	defer func() {
		if err := exec.Command("rm", "-rf", secondAssetPath).Run(); err != nil {
			t.Logf("failed to remove %s: %v", secondAssetPath, err)
		}
	}()

	input := env.WithRegistries(fmt.Sprintf(`
kind: Object
spec:
- image: docker.io/*username*/kilnbase-e2e-tests-build
- image: docker.io/*username*/kilnbase-e2e-tests-build2
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Sources
sources:
- image: docker.io/*username*/kilnbase-e2e-tests-build
  path: %s
  bazel:
    run:
      target: :simple-app
- image: docker.io/*username*/kilnbase-e2e-tests-build2
  path: %s
  bazel:
    run:
      target: :simple-app
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageDestinations
destinations:
- image: docker.io/*username*/kilnbase-e2e-tests-build
- image: docker.io/*username*/kilnbase-e2e-tests-build2
`, assetPath, secondAssetPath))

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	out = strings.Replace(out, regexp.MustCompile("sha256:[a-z0-9]{64}").FindString(out), "SHA256-REPLACED", -1)

