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
