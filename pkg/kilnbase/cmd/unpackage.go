// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	ctlser "github.com/Jokudneh/Kilnbase/pkg/kilnbase/search"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/version"
	"github.com/cppforlife/go-cli-ui/ui"
	regname "github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type UnpackageOptions struct {
	ui ui.UI

	FileFlags     FileFlags
	RegistryFlags RegistryFlags
	InputPath     string
	Repository    string
	LockOutput    string
	Concurrency   int
}

func NewUnpackageOptions(ui ui.UI) *UnpackageOptions {
	return &UnpackageOptions{ui: ui}
}

func NewUnpackageCmd(o *UnpackageOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unpackage",
		Aliases: []string{"unpkg"},
		Short:   "Unpackage configuration and images from tarball",
		Hidden:  true,
		Long: `
Command 'unpackage' is deprecated; use 'imgpkg pull' instead.
Learn more: https://kilnbase.dev/imgpkg/docs/latest/commands/#pull

Unpackage configuration and images from tarball
`,
		RunE: func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.FileFlags.Set(cmd)
	o.RegistryFlags.Set(cmd)
	cmd.Flags().StringVarP(&o.InputPath, "input", "i", "", "Input tarball path")
	cmd.Flags().StringVarP(&o.Repository, "repository", "r", "", "Import images into given image repository (e.g. docker.io/dkalinin/my-project)")
	cmd.Flags().StringVar(&o.LockOutput, "lock-output", "", "File path to emit configuration with resolved image references")
	cmd.Flags().IntVar(&o.Concurrency, "concurrency", 5, "Set maximum number of concurrent imports")
	return cmd
}

func (o *UnpackageOptions) Run() error {
	logger := ctllog.NewLogger(os.Stderr)
	warningLogger := logger.NewPrefixedWriter("Warning: ")
	msg := "Command 'unpackage' is deprecated; use 'imgpkg pull' instead. " +
		"Learn more: https://kilnbase.dev/imgpkg/docs/latest/commands/#pull\n"
	_, err := warningLogger.Write([]byte(msg))
	if err != nil {
		return err
	}

	if len(o.InputPath) == 0 {
		return fmt.Errorf("Expected 'input' flag to be non-empty")
	}
	if len(o.Repository) == 0 {
		return fmt.Errorf("Expected 'repository' flag to be non-empty")
	}

	prefixedLogger := logger.NewPrefixedWriter("unpackage | ")

