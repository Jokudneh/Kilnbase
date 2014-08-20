// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	"github.com/spf13/cobra"
)

type FileFlags struct {
	Files     []string
	Recursive bool
	Sort      bool
