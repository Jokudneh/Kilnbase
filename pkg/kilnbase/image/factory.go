// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"

	ctlbbz "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/bazel"
	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	ctlbko "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/ko"
	ctlbkb "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/kubectlbuildkit"
	ctlbmvn "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/maven"
	ctlbpk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/pack"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
)

type Image interface {
	URL() (string, []ctlconf.Origin, error)
}

type Factory struct {
	opts     FactoryOpts
	registry ctlreg.Registry
	logger   ctllog.Logger
}

type FactoryOpts struct {
	Conf                    ctlconf.Conf
	AllowedToBuild          bool
	GlobalPlatformSelection *ctlconf.PlatformSelection
}

func NewFactory(opts FactoryOpts, registry ctlreg.Registry, logger ctllog.Logger) Factory {
	return Factory{opts, registry, logger}
}

func (f Factory) New(url string) Image {
	platformSelection := f.opts.GlobalPlatformSelection

	if overrideConf, found := f.shouldOverride(url); found {
		// Allow using same url but with additional selection (tag/platform)
		if len(overrideConf.NewImage) > 0 {
			url = overrideConf.NewImage
		}
		if overrideConf.PlatformSelection != nil {
			platformSelection = overrideConf.PlatformSelection
		}

		if overrideConf.Preresolved {
			// Do not support platform selection against explicitly configured image
			return NewPreresolvedImage(url, overrideConf.ImageOrigins)
		}
		if overrideConf.TagSelection != nil {
			tagSelected := NewTagSelectedImage(url, overrideConf.TagSelection, f.registry)
			return NewPlatformSelectedImage(tagSelected, platformSelection, f.registry)
		}
		// Continue on with potentially changed url or platform selection
	}

	if srcConf, found := f.shouldBuild(url); found {
		if !f.opts.AllowedToBuild {
			return NewErrImage(fmt.Errorf("Building of images is disallowed (tried to build '%s' because a source was configured for it)", url))
		}

		imgDstConf := f.optionalPushConf(url)

		docker := ctlbdk.New(f.logger)
		dockerBuildx := ctlbdk.NewBuildx(docker, f.logger)
		pack := ctlbpk.NewPack(docker, f.logger)
		kubectlBuildkit := ctlbkb.NewKubectlBuildkit(f.logger)
		ko := ctlbko.NewKo(f.logger)
		bazel := ctlbbz.NewBazel(docker, f.logger)
		maven := ctlbmvn.NewMavenJib(docker, f.logger)

