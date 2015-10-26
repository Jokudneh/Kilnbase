// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	GitRepoRemoteURLUnknown = "<unknown>"
	GitRepoHeadSHANoCommits = "<no commits>"
)

type GitRepo struct {
	dirPath string // CWD when exec-ing to avoid --git-dir weirdness
}

func NewGitRepo(dirPath string) GitRepo {
	return GitRepo{dirPath}
}

func (r GitRepo) RemoteURL() (string, error) {
	// TODO support other remotes besides origin?
	stdout, stderr, err := r.runCmd([]string{"ls-remote", "--get-url"})
	if err != nil {
		// Same message is returned if it's not a git repo
		if r.IsValid() && strings.Contains(stderr, "No remote configured to list refs from") {
			return GitRepoRemoteURLUnknown, nil
		}
		return "", r.error("Determining remote: %s (stderr '%s')", err, stderr)
	}

	return strings.TrimSpace(stdout), nil
}

func (r GitRepo) HeadSHA() (string, error) {
	stdout, stderr, err := r.runCmd([]string{"rev-parse", "HEAD"})
	if err != nil {
		listStdout, _, listErr := r.runCmd([]string{"rev-list", "-n", "1", "--all"})
		if listErr == nil && len(strings.TrimSpace(listStdout)) == 0 {
			return GitRepoHeadSHANoCommits, nil
		}
		return "", r.error("Checking HEAD commit: %s (stderr '%s')", err, stderr)
	}

	return strings.TrimSpace(stdout), nil
}

