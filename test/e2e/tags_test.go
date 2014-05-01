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
