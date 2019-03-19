//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestAdditionalImageTags(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, logger}

	var published map[string]interface{}

	logger.Section("build and tag image", func() {
		input := env.WithRegistries(`
image: docker.io/*username*/kilnbase-e2e-tests-build
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Sources
sources:
- image: docker.io/*username*/kilnbase-e2e-tests-build
  path: assets/simple-app
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageDestinations
destinations:
- image: docker.io/*username*/kilnbase-e2e-tests-build
  tags:
  - staging
  - v1.0.0
`)

		out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
			StdinReader: strings.NewReader(input),
		})

		maskedOut := strings.Replace(out, regexp.MustCompile("sha256:[a-z0-9]{64}").FindString(out), "SHA256-REPLACED", -1)

		expectedOut := env.WithRegistries(`---
image: index.docker.io/*username*/kilnbase-e2e-tests-build@SHA256-REPLACED
