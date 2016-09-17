// Copyright 2026 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedactIPv6AndFragmentCredentials(t *testing.T) {
	inputs := []string{
		`Get "http://[::1]:8080/input?token=query-canary": failed`,
		`redirect [https://user:password-canary@[::1]/?q=query-canary]`,
		`Get "https://example.com/input#fragment-canary": failed`,
		`Get "http://[::1]:bad/input?token=query-canary": failed`,
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			require.NotContains(t, redactErrorMessage(input), "canary")
		})
	}
	require.Equal(t,
		`redirect [https://_redacted_@[::1]/?_redacted_]`,
		redactErrorMessage(inputs[1]))
}

func TestHTTPDescriptionOmitsFragmentCredentials(t *testing.T) {
	source := NewHTTPFileSource(
		"https://example.com/input#fragment-canary")
	require.Equal(t, "HTTP URL 'https://example.com/input'",
		source.Description())
}

