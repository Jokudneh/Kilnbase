//go:build e2e

// Copyright 2026 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestMavenJibBuildSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	assetPath := "assets/maven-app"
	secondAssetPath := "assets/maven-app-2"
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
  maven:
    run:
      target: .
- image: docker.io/*username*/kilnbase-e2e-tests-build2
  path: %s
  maven:
    run:
      target: .
`, assetPath, secondAssetPath))

	out, _ := kilnbase.RunWithOpts(
		[]string{"-f", "-", "--images-annotation=false"},
		RunOpts{StdinReader: strings.NewReader(input)},
	)

	reSHA := regexp.MustCompile("sha256\\-[a-z0-9]{64}")
	out = reSHA.ReplaceAllString(out, "SHA256-REPLACED")
	reKilnbase := regexp.MustCompile(
		"kilnbase:(.+)-kilnbase-e2e(\\-.*)-SHA256-REPLACED")
	out = reKilnbase.ReplaceAllString(out, "kilnbase:img-title-SHA256-REPLACED")

	expectedOut := `---
kind: Object
spec:
- image: kilnbase:img-title-SHA256-REPLACED
- image: kilnbase:img-title-SHA256-REPLACED
`

	if out != expectedOut {
		t.Fatalf("Expected >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestMavenJibBuildAndPushSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	assetPath := "assets/maven-app"
	secondAssetPath := "assets/maven-app-2"
	if err := exec.Command("cp", "-r",
