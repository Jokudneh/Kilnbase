// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/cmd"
	uierrs "github.com/cppforlife/go-cli-ui/errors"
	"github.com/cppforlife/go-cli-ui/ui"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	log.SetOutput(io.Discard)

	// TODO logs
	// TODO log flags used

	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()

	options := cmd.NewKilnbaseOptions(confUI)
	command := cmd.NewKilnbaseCmd(options)

	err := command.Execute()
	if err != nil {
		multiLineErr := uierrs.NewMultiLineError(err)
		if options.UIFlags.JSON {
			_, _ = fmt.Fprintf(os.Stderr, "kilnbase: Error: %s\n", multiLineErr)
		} else {
			confUI.ErrorLinef("kilnbase: Error: %s", multiLineErr)
		}
