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
}

func (i DescribedImage) Ref() string { return i.desc.Refs[0] }

// Layers returns the ordered collection of filesystem layers that comprise this image.
// The order of the list is oldest/base layer first, and most-recent/top layer last.
func (i DescribedImage) Layers() ([]regv1.Layer, error) {
	var layers []regv1.Layer
	for _, layerTD := range i.desc.Layers {
		var layer regv1.Layer
		if layerTD.IsDistributable() {
			layerFile, err := i.layerProvider.FindLayer(layerTD)
			if err != nil {
				return nil, err
			}
			layer = NewDescribedLayer(layerTD, layerFile)
		} else {
			layer = NewForeignDescribedLayer(layerTD)
		}
		layers = append(layers, layer)
	}
	return layers, nil
}

// MediaType of this image's manifest.
func (i DescribedImage) MediaType() (types.MediaType, error) {
	return types.MediaType(i.desc.Manifest.MediaType), nil
