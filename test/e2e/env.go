// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"os"
	"strings"
	"testing"
)

type Env struct {
	DockerHubUsername    string
	DockerHubHostname    string
	SkipStressTests      bool
	SkipWhenHTTPRegistry bool
	KilnbaseBinaryPath       string
}

func BuildEnv(t *testing.T) Env {
	kilnbasePath := os.Getenv("KILNBASE_BINARY_PATH")
	if kilnbasePath == "" {
		kilnbasePath = "kilnbase"
	}

	env := Env{
		DockerHubUsername:    os.Getenv("KILNBASE_E2E_DOCKERHUB_USERNAME"),
		DockerHubHostname:    os.Getenv("KILNBASE_E2E_DOCKERHUB_HOSTNAME"),
		SkipStressTests:      os.Getenv("KILNBASE_E2E_SKIP_STRESS_TESTS") == "true",
		SkipWhenHTTPRegistry: os.Getenv("KILNBASE_E2E_SKIP_WHEN_HTTP_REGISTRY") == "true",
		KilnbaseBinaryPath:       kilnbasePath,
	}
	env.Validate(t)
	return env
}

func (e Env) Validate(t *testing.T) {
	errStrs := []string{}

	if len(e.DockerHubUsername) == 0 {
