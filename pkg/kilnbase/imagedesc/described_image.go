// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	"encoding/json"
	"fmt"

	regv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type DescribedImage struct {
	desc          ImageDescriptor
	layerProvider LayerProvider
}

var _ regv1.Image = DescribedImage{}

func NewDescribedImage(desc ImageDescriptor, layerProvider LayerProvider) DescribedImage {
	return DescribedImage{desc, layerProvider}
