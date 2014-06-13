// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"io"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/version"
	"github.com/cppforlife/cobrautil"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
)

type KilnbaseOptions struct {
	ui      *ui.ConfUI
	UIFlags UIFlags
}

