// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/util"
	regname "github.com/google/go-containerregistry/pkg/name"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
	regtran "github.com/google/go-containerregistry/pkg/v1/remote/transport"
	regtypes "github.com/google/go-containerregistry/pkg/v1/types"
	"golang.org/x/sync/errgroup"
)

type Registry interface {
	Generic(regname.Reference) (regv1.Descriptor, error)
	Index(regname.Reference) (regv1.ImageIndex, error)
	Image(regname.Reference) (regv1.Image, error)
}

type ImageRefDescriptors struct {
	registry Registry

	descs []ImageOrImageIndexDescriptor

	imageLayersLock sync.Mutex
	imageLayers     map[ImageLayerDescriptor]regv1.Layer
}

func NewImageRefDescriptorsFromBytes(data []byte) (*ImageRefDescriptors, error) {
	var descs []ImageOrImageIndexDescriptor

	err := json.Unmarshal(data, &descs)
	if err != nil {
		return nil, err
	}

	return &ImageRefDescriptors{descs: descs}, nil
}

func NewImageRefDescriptors(refs []regname.Reference, registry Registry) (*ImageRefDescriptors, error) {
	registry = errRegistry{registry}

	imageRefDescs := &ImageRefDescriptors{
		registry:    registry,
		imageLayers: map[ImageLayerDescriptor]regv1.Layer{},
	}

	var imageRefDescsLock sync.Mutex
	var wg errgroup.Group
	buildThrottle := util.NewThrottle(10)

	for _, ref := range refs {
		ref := ref //copy

		wg.Go(func() error {
			buildThrottle.Take()
			defer buildThrottle.Done()

			regDesc, err := registry.Generic(ref)
			if err != nil {
				return err
			}

			var td ImageOrImageIndexDescriptor

			if imageRefDescs.isImageIndex(regDesc) {
				imgIndexTd, err := imageRefDescs.buildImageIndex(ref, regDesc)
				if err != nil {
					return err
				}
				td = ImageOrImageIndexDescriptor{ImageIndex: &imgIndexTd}
			} else {
				imgTd, err := imageRefDescs.buildImage(ref)
				if err != nil {
					return err
				}
				td = ImageOrImageIndexDescriptor{Image: &imgTd}
			}

			imageRefDescsLock.Lock()
			imageRefDescs.descs = append(imageRefDescs.descs, td)
			imageRefDescsLock.Unlock()

			return nil
		})
	}

	err := wg.Wait()

	return imageRefDescs, err
}

func (ids *ImageRefDescriptors) Descriptors() []ImageOrImageIndexDescriptor {
	return ids.descs
}

func (ids *ImageRefDescriptors) buildImageIndex(ref regname.Reference, regDesc regv1.Descriptor) (ImageIndexDescriptor, error) {
	td := ImageIndexDescriptor{
		Refs:      []string{ref.Name()},
		MediaType: string(regDesc.MediaType),
		Digest:    regDesc.Digest.String(),
	}

	imgIndex, err := ids.registry.Index(ref)
	if err != nil {
		return td, err
	}

	rawManifest, err := imgIndex.RawManifest()
	if err != nil {
		return td, err
	}

	td.Raw = string(rawManifest)

	imgIndexManifest, err := imgIndex.IndexManifest()
	if err != nil {
		return td, err
	}

	for _, manDesc := range imgIndexManifest.Manifests {
		if ids.isImageIndex(manDesc) {
			imgIndexTd, err := ids.buildImageIndex(ids.buildRef(ref, manDesc.Digest.String()), manDesc)
			if err != nil {
				return ImageIndexDescriptor{}, err
			}
			td.Indexes = append(td.Indexes, imgIndexTd)
		} else {
			imgTd, err := ids.buildImage(ids.buildRef(ref, manDesc.Digest.String()))
			if err != nil {
				return ImageIndexDescriptor{}, err
			}
			td.Images = append(td.Images, imgTd)
		}
	}

	return td, nil
}

func (ids *ImageRefDescriptors) buildImage(ref regname.Reference) (ImageDescriptor, error) {
	td := ImageDescriptor{}

	img, err := ids.registry.Image(ref)
	if err != nil {
		return td, err
	}

	cfgDigest, err := img.ConfigName()
	if err != nil {
		return td, err
	}
	cfgBlob, err := img.RawConfigFile()
	if err != nil {
		return td, err
	}

	manifestMediaType, err := img.MediaType()
	if err != nil {
		return td, err
	}
	manifestDigest, err := img.Digest()
	if err != nil {
		return td, err
	}
	manifestBlob, err := img.RawManifest()
	if err != nil {
		return td, err
	}

	td = ImageDescriptor{
		Refs: []string{ref.String()},

		Config: ConfigDescriptor{
			Digest: cfgDigest.String(),
			Raw:    string(cfgBlob),
		},

		Manifest: ManifestDescriptor{
			MediaType: string(manifestMediaType),
			Digest:    manifestDigest.String(),
			Raw:       string(manifestBlob),
		},
	}

	layers, err := img.Layers()
	if err != nil {
		return td, err
	}

	for _, layer := range layers {
		layerMediaType, err := layer.MediaType()
		if err != nil {
			return td, err
		}
		layerDigest, err := layer.Digest()
		if err != nil {
			return td, err
		}
		layerDiffID, err := layer.DiffID()
		if err != nil {
			return td, err
		}
		layerSize, err := layer.Size()
		if err != nil {
			return td, err
		}

		layerTD := ImageLayerDescriptor{
			MediaType: string(layerMediaType),
			Digest:    layerDigest.String(),
			DiffID:    layerDiffID.String(),
			Size:      layerSize,
		}

		td.Layers = append(td.Layers, layerTD)

		ids.imageLayersLock.Lock()
		ids.imageLayers[layerTD] = layer
		ids.imageLayersLock.Unlock()
	}

	return td, nil
}

func (ImageRefDescriptors) isImageIndex(regDesc regv1.Descriptor) bool {
	switch regDesc.MediaType {
	case regtypes.OCIImageIndex, regtypes.DockerManifestList:
		return true
	}
	return false
}

type wrappedCompressedLayerContents struct {
	layer regv1.Layer
}

var _ LayerContents = wrappedCompressedLayerContents{}

func (lc wrappedCompressedLayerContents) Open() (io.ReadCloser, error) {
	rc, err := lc.layer.Compressed()
	if err != nil {
		return nil, fmt.Errorf("Getting compressed layer: %s", err)
	}
	return rc, nil
}

func (ids *ImageRefDescriptors) FindLayer(layerTD ImageLayerDescriptor) (LayerContents, error) {
	ids.imageLayersLock.Lock()
	defer ids.imageLayersLock.Unlock()
