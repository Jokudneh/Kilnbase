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

func TestPackBuildAndPushSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	// Copy asset to avoid parallel build interference
	assetPath := "assets/simple-app"
	secondAssetPath := "assets/simple-app-2"
	copyCmd := exec.Command("cp", "-r", assetPath, secondAssetPath)
	if err := copyCmd.Run(); err != nil {
		t.Fatalf("failed to copy asset from %s to %s: %v",
			assetPath, secondAssetPath, err)
	}
	defer func() {
		if err := exec.Command("rm", "-rf", secondAssetPath).Run(); err != nil {
			t.Logf("failed to remove temporary asset path %s: %v",
				secondAssetPath, err)
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
