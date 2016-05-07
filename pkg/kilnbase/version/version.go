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
		// Version was set via ldflags, just return it.
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return defaultVersion
	}

	// Anything else.
	for _, dep := range info.Deps {
		if dep.Path == moduleName {
