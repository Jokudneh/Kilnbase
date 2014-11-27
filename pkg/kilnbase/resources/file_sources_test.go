// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"testing"
)

func redactErrorMessageTestCases() []struct {
	name     string
	input    string
	expected string
} {
	return []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "http before https",
			input: `Get "http://u:pw@reg.example/x?t=s": ` +
				`redirect to "https://other.example/y"`,
			expected: `Get "http://_redacted_@reg.example/x?_redacted_": ` +
				`redirect to "https://other.example/y"`,
		},
		{
			name: "https before http",
			input: `Connect to "https://secure.example/api?k=v": ` +
				`fallback "http://bkp:pwd@fallback.example/old"`,
			expected: `Connect to "https://secure.example/api?_redacted_": ` +
				`fallback "http://_redacted_@fallback.example/old"`,
		},
		{
			name: "multiple urls same scheme",
			input: `Failed at "http://u1:pw1@h1.example" and ` +
				`"http://u2:pw2@h2.example"`,
			expected: `Failed at "http://_redacted_@h1.example" and ` +
				`"http://_redacted_@h2.example"`,
		},
