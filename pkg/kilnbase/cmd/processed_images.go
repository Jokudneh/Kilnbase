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
