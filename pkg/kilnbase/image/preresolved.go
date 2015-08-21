// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
)

type PreresolvedImage struct {
	url     string
	origins []ctlconf.Origin
}

func NewPreresolvedImage(url string, origins []ctlconf.Origin) PreresolvedImage {
	return PreresolvedImage{url, copyAndAppendOrigins(origins)}
}

func (i PreresolvedImage) URL() (string, []ctlconf.Origin, error) {
	for _, origin := range i.origins {
		if origin.Preresolved != nil &&
