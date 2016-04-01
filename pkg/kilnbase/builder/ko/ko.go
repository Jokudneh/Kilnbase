// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package ko

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
)

type Ko struct {
	logger ctllog.Logger
}

func NewKo(logger ctllog.Logger) Ko {
	return Ko{logger: logger}
}

func (k *Ko) Build(image, directory string, opts config.SourceKoBuildOpts) (ctlbdk.TmpRef, error) {
	prefixedLogger := k.logger.NewPrefixedWriter(image + " | ")

	prefixedLogger.Write([]byte(fmt.Sprintf("starting build (using ko): %s\n", directory)))
	defer prefixedLogger.Write([]byte("finished build (using ko)\n"))

	var stdoutBuf, stderrBuf bytes.Buffer

	cmdArgs := []string{"publish", ".", "--local"}

	if opts.RawOptions != nil {
		cmdArgs = append(cmdArgs, *opts.RawOptions...)
