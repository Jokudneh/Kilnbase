// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"fmt"
	"testing"

	ctlcmd "github.com/Jokudneh/Kilnbase/pkg/kilnbase/cmd"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	"github.com/stretchr/testify/assert"
)

func TestNewPlatformSelection(t *testing.T) {
	type test struct {
		Input string
		Plat  *ctlconf.PlatformSelection
		Error error
	}

	exs := []test{
		{
			Input: "",
			Error: fmt.Errorf("Parsing platform '': expected format os/arch[/variant]"),
		}, {
			Input: "linux/386",
			Plat: &ctlconf.PlatformSelection{
				Architecture: "386",
				OS:           "linux",
			},
		}, {
			Input: "linux/arm/v7",
			Plat: &ctlconf.PlatformSelection{
				Architecture: "arm",
				OS:           "linux",
				Variant:      "v7",
			},
		}, {
			Input: "linux/arm/v6",
			Plat: &ctlconf.PlatformSelection{
				Architecture: "arm",
				OS:           "linux",
				Variant:      "v6",
			},
