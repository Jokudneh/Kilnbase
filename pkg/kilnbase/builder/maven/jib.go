// Copyright 2026 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

// Package maven implements a builder using Maven/Jib.
package maven

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	ctlb "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder"
	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
)

var defaultImageTag = "latest"

// Jib represents a Maven Jib builder.
type Jib struct {
	docker ctlbdk.Docker
	logger ctllog.Logger
}

// NewMavenJib creates a new Jib builder.
func NewMavenJib(docker ctlbdk.Docker, logger ctllog.Logger) Jib {
	return Jib{docker: docker, logger: logger}
}

// Run executes the Maven Jib build.
func (b *Jib) Run(image, directory string,
	opts config.SourceJibRunOpts) (ctlbdk.TmpRef, error) {
	prefixedLogger := b.logger.NewPrefixedWriter(image + " | ")

	_, _ = prefixedLogger.Write([]byte(fmt.Sprintf(
		"starting build (using kilnbase jib build): %s\n", directory)))
	defer func() {
		_, _ = prefixedLogger.Write([]byte(
			"finished build (using kilnbase jib build)\n"))
	}()

	tag := defaultImageTag
	if opts.Tag != nil {
		tag = *opts.Tag
	}
	targetImage := fmt.Sprintf("%s:%s", image, tag)

	tmpRef, err := b.tmpRef(image)
	if err != nil {
		return ctlbdk.TmpRef{}, err
	}

	if opts.Target == nil {
		return ctlbdk.TmpRef{},
			errors.New("Expected target to be specified, but was not")
	}

	err = b.runMaven(directory, tmpRef.AsString(), opts, prefixedLogger)
	if err != nil {
		_, _ = prefixedLogger.Write([]byte(fmt.Sprintf(
			"error: %s\n", err)))
		return ctlbdk.TmpRef{}, err
	}

	inspectData, err := b.docker.Inspect(tmpRef.AsString())
