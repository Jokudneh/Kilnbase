// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"kilnbase.dev/imgpkg/pkg/imgpkg/lockconfig"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlimg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/image"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	ctlser "github.com/Jokudneh/Kilnbase/pkg/kilnbase/search"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/version"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type ResolveOptions struct {
	ui ui.UI

	FileFlags         FileFlags
	RegistryFlags     RegistryFlags
	AllowedToBuild    bool
	BuildConcurrency  int
	ImagesAnnotation  bool
	OriginsAnnotation bool
	ImageMapFile      string
	LockOutput        string
	ImgpkgLockOutput  string
	UnresolvedInspect bool
	Platform          string
}

func NewResolveOptions(ui ui.UI) *ResolveOptions {
	return &ResolveOptions{ui: ui}
}

func NewResolveCmd(o *ResolveOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Build images and update references",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.FileFlags.Set(cmd)
	o.RegistryFlags.Set(cmd)
	cmd.Flags().BoolVar(&o.AllowedToBuild, "build", true, "Allow building of images")
	cmd.Flags().IntVar(&o.BuildConcurrency, "build-concurrency", 4, "Set maximum number of concurrent builds")
	cmd.Flags().BoolVar(&o.ImagesAnnotation, "images-annotation", true, "Annotate resources with images annotation")
	cmd.Flags().BoolVar(&o.OriginsAnnotation, "origins-annotation", true, "Include origins annotation")
	cmd.Flags().StringVar(&o.ImageMapFile, "image-map-file", "", "Set image map file (/cnab/app/relocation-mapping.json in CNAB)")
	cmd.Flags().StringVar(&o.LockOutput, "lock-output", "", "File path to emit configuration with resolved image references")
	cmd.Flags().StringVar(&o.ImgpkgLockOutput, "imgpkg-lock-output", "", "File path to emit images lockfile with resolved image references")
	cmd.Flags().BoolVar(&o.UnresolvedInspect, "unresolved-inspect", false, "List image references found in inputs")
	cmd.Flags().StringVar(&o.Platform, "platform", "", "Apply platform selection to image indexes")
	return cmd
}

func (o *ResolveOptions) Run() error {
	if o.ImgpkgLockOutput != "" && o.LockOutput != "" {
		msg := "Can only output one lockfile type; provide only one of " +
			"'--lock-output' or '--imgpkg-lock-output'"
		return errors.New(msg)
	}
	logger := ctllog.NewLogger(os.Stderr)
	prefixedLogger := logger.NewPrefixedWriter("resolve | ")

	resBss, err := o.ResolveResources(&logger, prefixedLogger)
	if err != nil {
		return err
	}

	// Print all resources as one YAML stream
	for _, resBs := range resBss {
		resBs = append([]byte("---\n"), resBs...)
		o.ui.PrintBlock(resBs)
	}

	return nil
}

func (o *ResolveOptions) ResolveResources(logger *ctllog.Logger, pLogger *ctllog.PrefixWriter) ([][]byte, error) {
	nonConfigRs, conf, err := o.FileFlags.ResourcesAndConfig()
	if err != nil {
		return nil, err
	}

	conf, err = o.withImageMapConf(conf)
	if err != nil {
		return nil, err
	}

	registry, err := ctlreg.NewRegistry(o.RegistryFlags.AsRegistryOpts())
	if err != nil {
		return nil, err
	}

	opts := ctlimg.FactoryOpts{
		Conf:           conf,
		AllowedToBuild: o.AllowedToBuild,
	}
	if len(o.Platform) > 0 {
		opts.GlobalPlatformSelection, err = NewPlatformSelection(o.Platform)
		if err != nil {
			return nil, err
		}
	}
	imgFactory := ctlimg.NewFactory(opts, registry, *logger)

	imageURLs, err := o.collectImageReferences(nonConfigRs, conf)
	if err != nil {
		return nil, err
	}

	if o.UnresolvedInspect {
		output, err := imageURLs.Bytes()
		if err != nil {
			return nil, err
		}
		o.ui.PrintBlock(output)
		return nil, nil
	}

	resolvedImages, err := o.resolveImages(imageURLs, imgFactory)
	if err != nil {
		return nil, err
	}

	// Record final image transformation
	for _, pair := range resolvedImages.All() {
		pLogger.WriteStr("final: %s -> %s\n", pair.UnprocessedImageURL.URL, pair.Image.URL)
