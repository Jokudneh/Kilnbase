// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

const (
	ImagesAnnKey = "kilnbase.k14s.io/images"
)

type ResourceWithImages struct {
	contents map[string]interface{}
	images   []Image
}

func NewResourceWithImages(contents map[string]interface{}, images []Image) ResourceWithImages {
	// sort images lexicographically (by URL) to avoid unnecessary annotation changes
	sort.Slice(images, func(i, j int) bool {
		return images[i].URL < images[j].URL
	})
	return ResourceWithImages{contents, images}
}
