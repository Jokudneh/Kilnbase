// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"path/filepath"

	ctlbbz "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/bazel"
	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	ctlbko "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/ko"
	ctlbkb "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/kubectlbuildkit"
	ctlbmvn "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/maven"
	ctlbpk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/pack"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
)

type BuiltImage struct {
	url         string
	buildSource ctlconf.Source
	imgDst      *ctlconf.ImageDestination

	docker          ctlbdk.Docker
	dockerBuildx    ctlbdk.Buildx
	pack            ctlbpk.Pack
	kubectlBuildkit ctlbkb.KubectlBuildkit
	ko              ctlbko.Ko
	bazel           ctlbbz.Bazel
	jib             ctlbmvn.Jib
}

// BuildersOpts contains all the builders used to construct a BuiltImage.
type BuildersOpts struct {
	Docker          ctlbdk.Docker
	DockerBuildx    ctlbdk.Buildx
	Pack            ctlbpk.Pack
	KubectlBuildkit ctlbkb.KubectlBuildkit
	Ko              ctlbko.Ko
	Bazel           ctlbbz.Bazel
	Jib             ctlbmvn.Jib
}

// NewBuiltImage creates a new BuiltImage.
func NewBuiltImage(url string, buildSource ctlconf.Source,
	imgDst *ctlconf.ImageDestination, builders BuildersOpts) BuiltImage {
	return BuiltImage{
		url:             url,
		buildSource:     buildSource,
		imgDst:          imgDst,
		docker:          builders.Docker,
		dockerBuildx:    builders.DockerBuildx,
		pack:            builders.Pack,
		kubectlBuildkit: builders.KubectlBuildkit,
		ko:              builders.Ko,
		bazel:           builders.Bazel,
		jib:             builders.Jib,
	}
}

func (i BuiltImage) URL() (string, []ctlconf.Origin, error) {
	origins, err := i.sources()
	if err != nil {
		return "", nil, err
	}

	urlRepo, _ := URLRepo(i.url)

	switch {
	case i.buildSource.Pack != nil:
		opts := ctlbpk.PackBuildOpts{
			Builder:    i.buildSource.Pack.Build.Builder,
			Buildpacks: i.buildSource.Pack.Build.Buildpacks,
			ClearCache: i.buildSource.Pack.Build.ClearCache,
			RawOptions: i.buildSource.Pack.Build.RawOptions,
		}

		dockerTmpRef, err := i.pack.Build(urlRepo, i.buildSource.Path, opts)
		if err != nil {
			return "", nil, err
		}

		return i.optionalPushWithDocker(dockerTmpRef, origins)

	case i.buildSource.KubectlBuildkit != nil:
		url, err := i.kubectlBuildkit.BuildAndPush(
			urlRepo, i.buildSource.Path, i.imgDst, *i.buildSource.KubectlBuildkit)
		return url, origins, err

	case i.buildSource.Ko != nil:
		dockerTmpRef, err := i.ko.Build(urlRepo, i.buildSource.Path, i.buildSource.Ko.Build)
		if err != nil {
			return "", nil, err
