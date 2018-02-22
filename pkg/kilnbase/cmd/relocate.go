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

type RelocateOptions struct {
	ui ui.UI

	FileFlags     FileFlags
	RegistryFlags RegistryFlags
	Repository    string
	LockOutput    string
	Concurrency   int
}

func NewRelocateOptions(ui ui.UI) *RelocateOptions {
	return &RelocateOptions{ui: ui}
}

func NewRelocateCmd(o *RelocateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "relocate",
		Long: `
Command 'relocate' is deprecated; use 'imgpkg copy' instead.
Learn more: https://kilnbase.dev/imgpkg/docs/latest/commands/#copy

Relocate images between two registries
`,
		Short:  "Relocate images between two registries",
		Hidden: true,
		RunE:   func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.FileFlags.Set(cmd)
	o.RegistryFlags.Set(cmd)
	cmd.Flags().StringVarP(&o.Repository, "repository", "r", "", "Import images into given image repository (e.g. docker.io/dkalinin/my-project)")
	cmd.Flags().StringVar(&o.LockOutput, "lock-output", "", "File path to emit configuration with resolved image references")
	cmd.Flags().IntVar(&o.Concurrency, "concurrency", 5, "Set maximum number of concurrent imports")
	return cmd
}

func (o *RelocateOptions) Run() error {
	logger := ctllog.NewLogger(os.Stderr)
	warningLogger := logger.NewPrefixedWriter("Warning: ")
	msg := "Command 'relocate' is deprecated; use 'imgpkg copy' instead. " +
		"Learn more: https://kilnbase.dev/imgpkg/docs/latest/commands/#copy\n"
	_, err := warningLogger.Write([]byte(msg))
	if err != nil {
		return err
	}

	// basic checks
	if len(o.Repository) == 0 {
		return fmt.Errorf("Expected repository flag to be non-empty")
	}
