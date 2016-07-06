// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
)

// ResolvedImage represents an image that will be resolved into url+digest
type ResolvedImage struct {
	url      string
	registry ctlreg.Registry
}

func NewResolvedImage(url string, registry ctlreg.Registry) ResolvedImage {
	return ResolvedImage{url, registry}
}

func (i ResolvedImage) URL() (string, []ctlconf.Origin, error) {
	tag, err := regname.NewTag(i.url, regname.WeakValidation)
	if err != nil {
