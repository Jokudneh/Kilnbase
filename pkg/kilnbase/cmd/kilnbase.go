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

func NewKilnbaseOptions(ui *ui.ConfUI) *KilnbaseOptions {
	return &KilnbaseOptions{ui: ui}
}

func NewDefaultKilnbaseCmd(ui *ui.ConfUI) *cobra.Command {
	return NewKilnbaseCmd(NewKilnbaseOptions(ui))
}

func NewKilnbaseCmd(o *KilnbaseOptions) *cobra.Command {
	cmd := NewResolveCmd(NewResolveOptions(o.ui))

	cmd.Use = "kilnbase"
	cmd.Short = "kilnbase prepares Docker images to deploy to Kubernetes"
	cmd.Version = version.Version

	// Affects children as well
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	// Disable docs header
	cmd.DisableAutoGenTag = true

	// We need to do a little work for the automatic "completion" command to function as designed.
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.SetOutput(uiBlockWriter{o.ui}) // setting output for cmd.Help()

	o.UIFlags.Set(cmd)

	cmd.AddCommand(NewInspectCmd(NewInspectOptions(o.ui)))
	cmd.AddCommand(NewPackageCmd(NewPackageOptions(o.ui)))
	cmd.AddCommand(NewUnpackageCmd(NewUnpackageOptions(o.ui)))
	cmd.AddCommand(NewVersionCmd(NewVersionOptions(o.ui)))
	cmd.AddCommand(NewRelocateCmd(NewRelocateOptions(o.ui)))

	// Last one runs first
	cobrautil.VisitCommands(cmd, cobrautil.ReconfigureCmdWithSubcmd)
	cobrautil.VisitCommands(cmd, cobrautil.DisallowExtraArgs)

	cobrautil.VisitCommands(cmd, cobrautil.WrapRunEForCmd(func(*cobra.Command, []string) error {
