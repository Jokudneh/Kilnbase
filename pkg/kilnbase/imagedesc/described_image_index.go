// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagedesc

import (
	"encoding/json"
	"fmt"

	regv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type DescribedImageIndex struct {
	desc    ImageIndexDescriptor
	images  []regv1.Image
	indexes []regv1.ImageIndex
}

