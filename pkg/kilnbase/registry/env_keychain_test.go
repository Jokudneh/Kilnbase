// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"testing"

	regname "github.com/google/go-containerregistry/pkg/name"
)

// Registry hostnames are case-insensitive DNS names (RFC 1035/1123). An
// operator-configured KILNBASE_REGISTRY_HOSTNAME_* value must resolve
// credentials for a target image reference even when the two differ only
// by case.
func TestEnvKeychainResolveIsCaseInsensitive(t *testing.T) {
	t.Setenv("KILNBASE_TEST_HOSTNAME", "MyRegistry.Example.com")
	t.Setenv("KILNBASE_TEST_USERNAME", "user")
	t.Setenv("KILNBASE_TEST_PASSWORD", "pass")

	k := NewEnvKeychain("KILNBASE_TEST")
