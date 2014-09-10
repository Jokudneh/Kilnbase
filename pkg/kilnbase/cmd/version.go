// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/version"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
)

type VersionOptions struct {
	ui ui.UI
}

func NewVersionOptions(ui ui.UI) *VersionOptions {
