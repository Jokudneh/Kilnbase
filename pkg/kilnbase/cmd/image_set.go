// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	ctlimg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/image"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagedesc"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/util"
	regname "github.com/google/go-containerregistry/pkg/name"
)

type ImageSet struct {
	concurrency int
	logger      *ctllog.PrefixWriter
}

func (o ImageSet) Relocate(foundImages *UnprocessedImageURLs,
	importRepo regname.Repository, registry ctlreg.Registry) (*ProcessedImages, error) {

	ids, err := o.Export(foundImages, registry)
	if err != nil {
		return nil, err
	}

	return o.Import(imagedesc.NewDescribedReader(ids, ids).Read(), importRepo, registry)
}

func (o ImageSet) Export(foundImages *UnprocessedImageURLs,
	registry ctlreg.Registry) (*imagedesc.ImageRefDescriptors, error) {

	o.logger.WriteStr("exporting %d images...\n", len(foundImages.All()))
	defer func() { o.logger.WriteStr("exported %d images\n", len(foundImages.All())) }()

	var refs []regname.Reference

	for _, img := range foundImages.All() {
		// Validate strictly as these refs were already resolved
		ref, err := regname.NewDigest(img.URL, regname.StrictValidation)
		if err != nil {
			return nil, err
		}

		o.logger.Write([]byte(fmt.Sprintf("will export %s\n", img.URL)))
		refs = append(refs, ref)
	}

	ids, err := imagedesc.NewImageRefDescriptors(refs, registry)
	if err != nil {
		return nil, fmt.Errorf("Collecting packaging metadata: %s", err)
	}

	return ids, nil
}

func (o *ImageSet) Import(imgOrIndexes []imagedesc.ImageOrIndex,
	importRepo regname.Repository, registry ctlreg.Registry) (*ProcessedImages, error) {

	importedImages := NewProcessedImages()

	o.logger.WriteStr("importing %d images...\n", len(imgOrIndexes))
	defer func() { o.logger.WriteStr("imported %d images\n", len(importedImages.All())) }()

	errCh := make(chan error, len(imgOrIndexes))
