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
	docker ctlbdk.Docker
	logger ctllog.Logger
}

func NewBazel(docker ctlbdk.Docker, logger ctllog.Logger) Bazel {
	return Bazel{docker: docker, logger: logger}
}

func (b *Bazel) Run(image, directory string, opts config.SourceBazelRunOpts) (ctlbdk.TmpRef, error) {
	prefixedLogger := b.logger.NewPrefixedWriter(image + " | ")

	prefixedLogger.Write([]byte(fmt.Sprintf("starting build (using bazel): %s\n", directory)))
	defer prefixedLogger.Write([]byte("finished build (using bazel)\n"))

	var imageID string
	{
		var stdoutBuf, stderrBuf bytes.Buffer

		cmdArgs := []string{"run"}

		if opts.Target == nil {
			return ctlbdk.TmpRef{}, fmt.Errorf("Expected target to be specified, but was not")
		}

		cmdArgs = append(cmdArgs, *opts.Target)

		if opts.RawOptions != nil {
			cmdArgs = append(cmdArgs, *opts.RawOptions...)
		}

		cmd := exec.Command("bazel", cmdArgs...)
		cmd.Dir = directory
		cmd.Stdout = io.MultiWriter(&stdoutBuf, prefixedLogger)
		cmd.Stderr = io.MultiWriter(&stderrBuf, prefixedLogger)

		err := cmd.Run()
		if err != nil {
