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
}

func (b *ImageQueue) Run(unprocessedImageURLs *UnprocessedImageURLs, numWorkers int) (*ProcessedImages, error) {
	b.outputImages = NewProcessedImages()
	b.outputErrs = nil

	queueCh := make(chan UnprocessedImageURL, numWorkers)
	workWg := sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		go b.worker(&workWg, queueCh)
	}

	for _, unprocessedImageURL := range unprocessedImageURLs.All() {
		workWg.Add(1)
		queueCh <- unprocessedImageURL
	}

	workWg.Wait()
	close(queueCh)

	return b.outputImages, errFromErrs(b.outputErrs)
}
