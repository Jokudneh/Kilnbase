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

	return strings.TrimSpace(string(yamlBytes))
}

type imageStruct struct {
	URL     string        `json:"url"`
	Origins []interface{} `json:"origins,omitempty"`
}

func (st imageStruct) equal(other imageStruct) bool {
	return st.URL == other.URL && reflect.DeepEqual(st.Origins, other.Origins)
}

func contains(structs []imageStruct, st imageStruct) bool {
	for _, other := range structs {
		if st.equal(other) {
			return true
		}
	}
