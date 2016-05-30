// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagetar"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
)

type TarImageSet struct {
	imageSet    ImageSet
	concurrency int
	logger      *ctllog.PrefixWriter
}

func (o TarImageSet) Export(foundImages *UnprocessedImageURLs,
	outputPath string, registry ctlreg.Registry) error {

	ids, err := o.imageSet.Export(foundImages, registry)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Creating file '%s': %s", outputPath, err)
	}

	err = outputFile.Close()
	if err != nil {
		return err
	}

