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
