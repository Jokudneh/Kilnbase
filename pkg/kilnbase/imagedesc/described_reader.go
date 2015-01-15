// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	regv1 "github.com/google/go-containerregistry/pkg/v1"
)

type DescribedReader struct {
	ids           *ImageRefDescriptors
	layerProvider LayerProvider
}

func NewDescribedReader(ids *ImageRefDescriptors, layerProvider LayerProvider) DescribedReader {
	return DescribedReader{ids, layerProvider}
}

func (r DescribedReader) Read() []ImageOrIndex {
	var result []ImageOrIndex
