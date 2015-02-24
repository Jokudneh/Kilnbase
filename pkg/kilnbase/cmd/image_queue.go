// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"sync"

	ctlimg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/image"
)

type ImageQueue struct {
	imgFactory ctlimg.Factory

	outputImages     *ProcessedImages
	outputImagesLock sync.Mutex

	outputErrs     []error
	outputErrsLock sync.Mutex
}

func NewImageQueue(imgFactory ctlimg.Factory) *ImageQueue {
	return &ImageQueue{imgFactory: imgFactory}
