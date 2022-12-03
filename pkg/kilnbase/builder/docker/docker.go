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
		if opts.Pull != nil && *opts.Pull {
			cmdArgs = append(cmdArgs, "--pull")
		}
		if opts.NoCache != nil && *opts.NoCache {
			cmdArgs = append(cmdArgs, "--no-cache")
		}
		if opts.File != nil {
			// Since docker command is executed with cwd of directory,
			// Dockerfile path doesnt need to be joined with it
			cmdArgs = append(cmdArgs, "--file", *opts.File)
		}
		if opts.RawOptions != nil {
			cmdArgs = append(cmdArgs, *opts.RawOptions...)
		}

		cmdArgs = append(cmdArgs, "--tag", tmpRef.AsString(), ".")

		cmd := exec.Command("docker", cmdArgs...)
		cmd.Dir = directory
		cmd.Stdout = io.MultiWriter(&stdoutBuf, prefixedLogger)
		cmd.Stderr = io.MultiWriter(&stderrBuf, prefixedLogger)

		if opts.Buildkit != nil {
			cmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")
		}

		err := cmd.Run()
		if err != nil {
			prefixedLogger.Write([]byte(fmt.Sprintf("error: %s\n", err)))
			return TmpRef{}, err
		}
	}

	inspectData, err := d.Inspect(tmpRef.AsString())
	if err != nil {
		prefixedLogger.Write([]byte(fmt.Sprintf("inspect error: %s\n", err)))
		return TmpRef{}, err
	}

	return d.RetagStable(tmpRef, image, inspectData.ID, prefixedLogger)
}

func (d Docker) RetagStable(tmpRef TmpRef, image, imageID string,
	prefixedLogger *ctllog.PrefixWriter) (TmpRef, error) {

	tb := ctlb.TagBuilder{}

	// Retag image with its sha256 to produce exact image ref if nothing has changed.
	// Seems that Docker doesn't like `kilnbase@sha256:...` format for local images.
	// Image hint at the beginning for easier sorting.
	stableTmpRef := TmpRef{"kilnbase:" + tb.CheckTagLen128(fmt.Sprintf(
		"%s-%s",
		tb.TrimStr(tb.CleanStr(image), 50),
		tb.CheckLen(tb.CleanStr(imageID), 72),
	))}

	{
		var stdoutBuf, stderrBuf bytes.Buffer

		cmd := exec.Command("docker", "tag", tmpRef.AsString(), stableTmpRef.AsString())
		cmd.Stdout = io.MultiWriter(&stdoutBuf, prefixedLogger)
		cmd.Stderr = io.MultiWriter(&stderrBuf, prefixedLogger)

		err := cmd.Run()
		if err != nil {
			prefixedLogger.Write([]byte(fmt.Sprintf("tag error: %s\n", err)))
			return TmpRef{}, err
		}
	}

	// Remove temporary tag to be nice to `docker images` output.
	// (No point in "untagging" digest reference)
	if !strings.HasPrefix(tmpRef.AsString(), "sha256:") {
		var stdoutBuf, stderrBuf bytes.Buffer

		cmd := exec.Command("docker", "rmi", tmpRef.AsString())
		cmd.Stdout = io.MultiWriter(&stdoutBuf, prefixedLogger)
		cmd.Stderr = io.MultiWriter(&stderrBuf, prefixedLogger)

		err := cmd.Run()
		if err != nil {
			prefixedLogger.Write([]byte(fmt.Sprintf("untag error: %s\n", err)))
			return TmpRef{}, err
		}
	}

	return stableTmpRef, nil
}

func (d Docker) Push(tmpRef TmpRef, imageDst string) (ImageDigest, error) {
	prefixedLogger := d.logger.NewPrefixedWriter(imageDst + " | ")

	tb := ctlb.TagBuilder{}

	// Generate random tag for pushed image.
	// TODO we are technically polluting registry with new tags.
	// Unfortunately we do not know digest upfront so cannot use kilnbase-sha256-... format.
	imageDstTagged, err := regname.NewTag(imageDst, regname.WeakValidation)
	badNameErr := (*regname.ErrBadName)(nil)
	if err == nil {
		randSuffix, err := tb.RandomStr50()
		if err != nil {
			return ImageDigest{}, fmt.Errorf("Generating image dst suffix: %s", err)
		}

		imageDstTag := fmt.Sprintf("kilnbase-%s", randSuffix)

		imageDstTagged, err = regname.NewTag(imageDst+":"+imageDstTag, regname.WeakValidation)
		if err != nil {
			return ImageDigest{}, fmt.Errorf(
				"Generating image dst tag '%s': %s", imageDst, err)
		}
	} else if errors.As(err, &badNameErr) {
		imageDstTagged, err = regname.NewTag(LowerCaseRegistry(imageDst))
		if err != nil {
			newError := fmt.Errorf(
				"Lower casing repository '%s' still failed: %w", imageDst, err)
			return ImageDigest{}, newError
		}
	}

	imageDst = imageDstTagged.Name()

	prefixedLogger.Write([]byte(fmt.Sprintf("starting push (using Docker): %s -> %s\n", tmpRef.AsString(), imageDst)))
	defer prefixedLogger.Write([]byte("finished push (using Docker)\n"))

	prevInspectData, err := d.Inspect(tmpRef.AsString())
	if err != nil {
		prefixedLogger.Write([]byte(fmt.Sprintf("inspect error: %s\n", err)))
		return ImageDigest{}, err
	}

	{
		var stdoutBuf, stderrBuf bytes.Buffer

		cmd := exec.Command("docker", "tag", tmpRef.AsString(), imageDst)
		cmd.Stdout = io.MultiWriter(&stdoutBuf, prefixedLogger)
		cmd.Stderr = io.MultiWriter(&stderrBuf, prefixedLogger)

		err := cmd.Run()
		if err != nil {
			prefixedLogger.Write([]byte(fmt.Sprintf("tag error: %s\n", err)))
			return ImageDigest{}, err
