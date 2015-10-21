// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	"fmt"
	"io"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imageutils/gzip"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imageutils/verify"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type DescribedLayer struct {
	desc     ImageLayerDescriptor
	contents LayerContents
}

var _ regv1.Layer = DescribedLayer{}

func NewDescribedLayer(desc ImageLayerDescriptor, contents LayerContents) DescribedLayer {
	return DescribedLayer{desc, contents}
}

