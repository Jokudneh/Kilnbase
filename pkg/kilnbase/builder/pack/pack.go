// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"

	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
)

var (
	// Example output that includes final digest:
	// [exporter] *** Images:
	// [exporter]       index.docker.io/library/myapp:latest - succeeded
	// [exporter] *** Image ID: 2be602fbc1ecffdf9cc1c8ccb1f1cd6fb1d0a2e76dccbfcc34898bf35c836beb
	// Image ID is printed here: https://github.com/buildpack/lifecycle/blob/4e449525af56096f7cf8a521900bf6216467f0d7/save.go#L39
	packImageID = regexp.MustCompile("Image ID: (sha256:)?([0-9a-z]+)")
)

type Pack struct {
	docker ctlbdk.Docker
	logger ctllog.Logger
}

type PackBuildOpts struct {
	Builder    *string
	Buildpacks *[]string
	ClearCache *bool
	RawOptions *[]string // pack build -h
}

func NewPack(docker ctlbdk.Docker, logger ctllog.Logger) Pack {
