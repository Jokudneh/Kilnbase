// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"reflect"
	"strings"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	"sigs.k8s.io/yaml"
)

type Images []Image

type Image struct {
	URL        string
	Origins    []ctlconf.Origin // empty when deserialized
	originsRaw []interface{}    // populated when deserialized
}

func (imgs Images) ForImage(url string) (Image, bool) {
	for _, img := range imgs {
		if img.URL == url {
			return img, true
		}
	}
	return Image{}, false
}

// TODO only works after deserialization
func (i Image) Description() string {
	yamlBytes, err := yaml.Marshal(i.originsRaw)
	if err != nil {
		return "[]" // TODO deal better?
	}
