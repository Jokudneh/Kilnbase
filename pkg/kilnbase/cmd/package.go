// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagedesc"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	ctlser "github.com/Jokudneh/Kilnbase/pkg/kilnbase/search"
	"github.com/cppforlife/go-cli-ui/ui"
	regname "github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
)

type PackageOptions struct {
	ui ui.UI

	FileFlags     FileFlags
	RegistryFlags RegistryFlags
	OutputPath    string
	Concurrency   int
}

var _ imagedesc.Registry = ctlreg.Registry{}

func NewPackageOptions(ui ui.UI) *PackageOptions {
	return &PackageOptions{ui: ui}
}

func NewPackageCmd(o *PackageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package",
		Aliases: []string{"pkg"},
		Short:   "Package configuration and images into tarball",
		Hidden:  true,
		Long: `
Command 'package' is deprecated; use 'imgpkg push' instead.
Learn more: https://kilnbase.dev/imgpkg/docs/latest/commands/#push

Package configuration and images into tarball
`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.FileFlags.Set(cmd)
	o.RegistryFlags.Set(cmd)
	cmd.Flags().StringVarP(&o.OutputPath, "output", "o", "", "Output tarball path")
	cmd.Flags().IntVar(&o.Concurrency, "concurrency", 5, "Set maximum number of concurrent imports")
	return cmd
}

func (o *PackageOptions) Run() error {
	logger := ctllog.NewLogger(os.Stderr)
	warningLogger := logger.NewPrefixedWriter("Warning: ")
	msg := "Command 'package' is deprecated; use 'imgpkg push' instead. " +
		"Learn more: https://kilnbase.dev/imgpkg/docs/latest/commands/#push\n"
