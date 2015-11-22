// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	"io"

	regv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type ForeignDescribedLayer struct {
	desc ImageLayerDescriptor
}

var _ regv1.Layer = ForeignDescribedLayer{}

func NewForeignDescribedLayer(desc ImageLayerDescriptor) ForeignDescribedLayer {
	return ForeignDescribedLayer{desc}
}

