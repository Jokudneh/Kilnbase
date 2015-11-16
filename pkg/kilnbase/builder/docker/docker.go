// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	ctlb "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	regname "github.com/google/go-containerregistry/pkg/name"
)

type Docker struct {
	logger ctllog.Logger
}

type BuildOpts struct {
	// https://docs.docker.com/engine/reference/commandline/build/
	Target     *string
	Pull       *bool
	NoCache    *bool
	File       *string
	Buildkit   *bool
	RawOptions *[]string
}

type TmpRef struct {
	val string
}

func NewTmpRef(val string) TmpRef {
	return TmpRef{val}
}

func (r TmpRef) AsString() string { return r.val }

type ImageDigest struct {
	val string
}

func (r ImageDigest) AsString() string { return r.val }

func New(logger ctllog.Logger) Docker {
	return Docker{logger}
}

func (d Docker) Build(image, directory string, opts BuildOpts) (TmpRef, error) {
	err := d.ensureDirectory(directory)
	if err != nil {
		return TmpRef{}, err
	}

	tb := ctlb.TagBuilder{}

	randPrefix50, err := tb.RandomStr50()
	if err != nil {
		return TmpRef{}, fmt.Errorf("Generating tmp image suffix: %s", err)
	}

	tmpRef := TmpRef{"kilnbase:" + tb.CheckTagLen128(fmt.Sprintf(
		"%s-%s",
		randPrefix50,
		tb.TrimStr(tb.CleanStr(image), 50),
	))}

	prefixedLogger := d.logger.NewPrefixedWriter(image + " | ")

	prefixedLogger.Write([]byte(fmt.Sprintf("starting build (using Docker): %s -> %s\n", directory, tmpRef.AsString())))
	defer prefixedLogger.Write([]byte("finished build (using Docker)\n"))

	{
		var stdoutBuf, stderrBuf bytes.Buffer

		cmdArgs := []string{"build"}

		if opts.Target != nil {
			cmdArgs = append(cmdArgs, "--target", *opts.Target)
		}
