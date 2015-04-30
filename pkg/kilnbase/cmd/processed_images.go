// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"sort"
	"sync"
)

type ProcessedImages struct {
	imgs     map[UnprocessedImageURL]Image
	imgsLock sync.Mutex
}

type ProcessedImageItem struct {
	UnprocessedImageURL
	Image
}

func NewProcessedImages() *ProcessedImages {
	return &ProcessedImages{imgs: map[UnprocessedImageURL]Image{}}
}

func (i *ProcessedImages) Add(unprocessedImageURL UnprocessedImageURL, img Image) {
	i.imgsLock.Lock()
	defer i.imgsLock.Unlock()

	i.imgs[unprocessedImageURL] = img
}

func (i *ProcessedImages) FindByURL(unprocessedImageURL UnprocessedImageURL) (Image, bool) {
	i.imgsLock.Lock()
