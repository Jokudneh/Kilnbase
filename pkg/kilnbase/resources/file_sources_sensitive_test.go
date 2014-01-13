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
