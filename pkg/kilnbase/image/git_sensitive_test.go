// Copyright 2026 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image_test

import (
	"strings"
	"testing"

	ctlimg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/image"
)

func TestGitProvenanceOmitsQueryCredentials(t *testing.T) {
	inputs := []string{
		"https://user:password-canary@example.com/repo" +
			"?token=query-canary#fragment-canary",
		"https://example.com/repo?token=query-canary",
		"git@example.com:repo?token=query-canary",
	}
	for _, raw := range inputs {
		if strings.Contains(ctlimg.GitRedactedRemoteURL(raw), "canary") {
			t.Fatal("provenance contains credentials")
		}
	}
