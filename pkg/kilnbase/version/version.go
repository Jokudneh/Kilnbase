// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"runtime/debug"
	"strings"
)

var (
	// Version can be set via:
	// -ldflags="-X 'github.com/Jokudneh/Kilnbase/pkg/kilnbase/version.Version=$TAG'"
	defaultVersion = "develop"
	Version        = ""
	moduleName     = "github.com/Jokudneh/Kilnbase"
)

func init() {
	Version = version()
}

func version() string {
	if Version != "" {
