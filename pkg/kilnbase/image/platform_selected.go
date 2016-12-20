// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
	regtypes "github.com/google/go-containerregistry/pkg/v1/types"
)

// PlatformSelectedImage selects specific image matching arch/platform
type PlatformSelectedImage struct {
	image     Image
	selection *ctlconf.PlatformSelection
	registry  ctlreg.Registry
}

func NewPlatformSelectedImage(image Image,
	selection *ctlconf.PlatformSelection, registry ctlreg.Registry) PlatformSelectedImage {

	return PlatformSelectedImage{image, selection, registry}
}

func (i PlatformSelectedImage) URL() (string, []ctlconf.Origin, error) {
	url, origins, err := i.image.URL()
	if err != nil {
		return url, origins, err
	}

	// Do nothing if platform selection is not requested
	if i.selection == nil {
		return url, origins, err
	}

	ref, err := regname.ParseReference(url)
	if err != nil {
		return "", nil, err
	}

	desc, err := i.registry.Generic(ref)
	if err != nil {
		return "", nil, err
	}

	switch desc.MediaType {
	case regtypes.OCIImageIndex, regtypes.DockerManifestList:
		imgIndex, err := i.registry.Index(ref)
		if err != nil {
			return "", nil, err
		}
