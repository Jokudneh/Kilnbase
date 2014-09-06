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
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
)

var (
	// Loaded image ID: sha256:328b5f47550c85cea5284911ad4d284ce20e8240d61d2610eb6cb4aa8b43c19e
	// Tagging 328b5f47550c85cea5284911ad4d284ce20e8240d61d2610eb6cb4aa8b43c19e as bazel:simple-app
	bazelImageID = regexp.MustCompile("Loaded image ID: (sha256:)([0-9a-z]+)")
)

type Bazel struct {
