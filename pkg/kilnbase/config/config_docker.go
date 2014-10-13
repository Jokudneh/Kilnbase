// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package config

type SourceDockerOpts struct {
	Build  SourceDockerBuildOpts
	Buildx *SourceDockerBuildxOpts
}

type SourceDockerBuildOpts struct {
