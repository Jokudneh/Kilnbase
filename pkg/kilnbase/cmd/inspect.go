// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	ctlser "github.com/Jokudneh/Kilnbase/pkg/kilnbase/search"
	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/spf13/cobra"
)

type InspectOptions struct {
	ui ui.UI

	FileFlags     FileFlags
	RegistryFlags RegistryFlags
}

func NewInspectOptions(ui ui.UI) *InspectOptions {
	return &InspectOptions{ui: ui}
}

func NewInspectCmd(o *InspectOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "inspect",
		Aliases: []string{"i", "is", "insp"},
		Short:   "Inspect images",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.FileFlags.Set(cmd)
	o.RegistryFlags.Set(cmd)
	return cmd
}

func (o *InspectOptions) Run() error {
	rs, conf, err := o.FileFlags.ResourcesAndConfig()
	if err != nil {
		return err
	}

	foundImages, err := o.findImages(rs, conf)
	if err != nil {
		return err
	}

	table := uitable.Table{
		Title:   "Images",
		Content: "images",

		Header: []uitable.Header{
			uitable.NewHeader("Image"),
			uitable.NewHeader("Metadata"),
			uitable.NewHeader("Resource"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 2, Asc: true},
