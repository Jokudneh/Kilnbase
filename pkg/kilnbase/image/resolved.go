// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
)

// ResolvedImage represents an image that will be resolved into url+digest
type ResolvedImage struct {
	url      string
	registry ctlreg.Registry
}

func NewResolvedImage(url string, registry ctlreg.Registry) ResolvedImage {
	return ResolvedImage{url, registry}
}

func (i ResolvedImage) URL() (string, []ctlconf.Origin, error) {
	tag, err := regname.NewTag(i.url, regname.WeakValidation)
	if err != nil {
		return "", nil, err
	}

	imgDescriptor, err := i.registry.Generic(tag)
	if err != nil {
		return "", nil, err
	}

	// Resolve image second time because some older registry can
	// return "random" digests that change for every request.
	// See https://github.com/Jokudneh/Kilnbase/issues/21 for details.
	imgDescriptor2, err := i.registry.Generic(tag)
	if err != nil {
		return "", nil, err
	}

	if imgDescriptor.Digest.String() != imgDescriptor2.Digest.String() {
		return "", nil, fmt.Errorf("Expected digest resolution to be consistent over two separate requests")
	}

	url, origins, err := NewDigestedImageFromParts(tag.Repository.String(), imgDescriptor.Digest.String()).URL()
	if err != nil {
		return "", nil, err
	}

	origins = append(origins, ctlconf.Origin{Resolved: &ctlconf.OriginResolved{URL: i.url, Tag: tag.TagStr()}})

