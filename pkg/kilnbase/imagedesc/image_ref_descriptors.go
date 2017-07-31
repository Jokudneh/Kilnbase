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
