// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"testing"
)

func TestLowerCaseRepository(t *testing.T) {
	for _, tc := range lowerCaseRepositoryCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LowerCaseRegistry(tc.input)
			if result != tc.expected {
				t.Errorf(
					"LowerCaseRegistry(%q) = %q, expected %q",
					tc.input, result, tc.expected)
			}
		})
	}
}

type lowerCaseRepositoryTestCase struct {
	name     string
	input    string
	expected string
}
