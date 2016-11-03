//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	out, _ := kilnbase.RunWithOpts([]string{"version"}, RunOpts{})

