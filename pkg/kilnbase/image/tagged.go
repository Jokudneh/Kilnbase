// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
)

// TaggedImage represents an image that will be tagged when its URL is requested
type TaggedImage struct {
	image    Image
	imgDst   ctlconf.ImageDestination
	registry ctlreg.Registry
}

func NewTaggedImage(image Image, imgDst ctlconf.ImageDestination, registry ctlreg.Registry) TaggedImage {
	return TaggedImage{image, imgDst, registry}
