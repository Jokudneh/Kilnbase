// Copyright 2026 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"strings"
	"testing"
)

func TestInvalidRegistryHostnameDoesNotExposeCredentials(t *testing.T) {
	for _, suffix := range []string{"", "_0"} {
		t.Run("hostname"+suffix, func(t *testing.T) {
			raw := "https://user:password-canary@example.com:bad" +
