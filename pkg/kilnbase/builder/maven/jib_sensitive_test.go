// Copyright 2026 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package maven

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	"github.com/stretchr/testify/require"
)

func TestMavenBannerDoesNotExposeArguments(t *testing.T) {
	const executableMode = 0700
	dir := t.TempDir()
	script := "#!/bin/sh\nprintf '%s\n' \"$@\" > \"$MAVEN_ARGS_FILE\"\n"
	err := os.WriteFile(filepath.Join(dir, "mvn"), []byte(script),
		executableMode)
	require.NoError(t, err)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	argsFile := filepath.Join(dir, "args")
	t.Setenv("MAVEN_ARGS_FILE", argsFile)
	opts := []string{
		"-Dpassword=password-canary",
		"-Durl=https://user:query-canary@example.com",
	}
	target := "."
	var log bytes.Buffer
	runOpts := config.SourceJibRunOpts{Target: &target, RawOptions: &opts}
	writer := logger.NewLogger(&log).NewPrefixedWriter("")
	err = (&Jib{}).runMaven(dir, "example.com/image", runOpts, writer)
	require.NoError(t, err)
