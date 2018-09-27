// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	"io"

	regv1 "github.com/google/go-containerregistry/pkg/v1"
	regv1types "github.com/google/go-containerregistry/pkg/v1/types"
)

type ImageOrIndex struct {
	Image *ImageWithRef
	Index *ImageIndexWithRef
}

type ImageWithRef interface {
	regv1.Image
	Ref() string
}

type ImageIndexWithRef interface {
	regv1.ImageIndex
	Ref() string
}

type LayerProvider interface {
	FindLayer(ImageLayerDescriptor) (LayerContents, error)
}

type LayerContents interface {
	Open() (io.ReadCloser, error)
}

type ImageOrImageIndexDescriptor struct {
	ImageIndex *ImageIndexDescriptor
	Image      *ImageDescriptor
}

type ImageIndexDescriptor struct {
	Refs    []string
	Images  []ImageDescriptor
	Indexes []ImageIndexDescriptor

	MediaType string
	Digest    string
	Raw       string
}

type ImageDescriptor struct {
	Refs   []string
	Layers []ImageLayerDescriptor

	Config   ConfigDescriptor
	Manifest ManifestDescriptor
}

type ImageLayerDescriptor struct {
	MediaType string
