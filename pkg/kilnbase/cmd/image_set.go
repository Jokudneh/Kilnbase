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
	importThrottle := util.NewThrottle(o.concurrency)

	for _, item := range imgOrIndexes {
		item := item // copy

		go func() {
			importThrottle.Take()
			defer importThrottle.Done()

			existingRef, err := regname.NewDigest(item.Ref())
			if err != nil {
				errCh <- err
				return
			}

			importDigestRef, err := o.importImage(item, existingRef, importRepo, registry)
			if err != nil {
				errCh <- fmt.Errorf("Importing image %s: %s", existingRef.Name(), err)
				return
			}

			importedImages.Add(UnprocessedImageURL{existingRef.Name()}, Image{URL: importDigestRef.Name()})
			errCh <- nil
		}()
	}

	for i := 0; i < len(imgOrIndexes); i++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}
	}

	return importedImages, nil
}

func (o *ImageSet) importImage(item imagedesc.ImageOrIndex,
	existingRef regname.Digest, importRepo regname.Repository,
	registry ctlreg.Registry) (regname.Digest, error) {

	itemDigest, err := item.Digest()
	if err != nil {
		return regname.Digest{}, err
	}

	importDigestRef, err := regname.NewDigest(fmt.Sprintf("%s@%s", importRepo.Name(), itemDigest))
	if err != nil {
		return regname.Digest{}, fmt.Errorf("Building new digest image ref: %s", err)
	}

	// Seems like AWS ECR doesnt like using digests for manifest uploads
	uploadTagRef, err := regname.NewTag(fmt.Sprintf("%s:kilnbase-%s-%s", importRepo.Name(), itemDigest.Algorithm, itemDigest.Hex))
	if err != nil {
		return regname.Digest{}, fmt.Errorf("Building upload tag image ref: %s", err)
	}

	o.logger.Write([]byte(fmt.Sprintf("importing %s -> %s...\n", existingRef.Name(), importDigestRef.Name())))
