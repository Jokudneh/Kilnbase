// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"
	"strings"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	regname "github.com/google/go-containerregistry/pkg/name"
)

const digestSep = "@"

type DigestedImage struct {
	nameWithDigest regname.Digest
	parseErr       error
}

var _ Image = DigestedImage{}

func MaybeNewDigestedImage(url string) *DigestedImage {
	nameWithDigest, err := regname.NewDigest(url, regname.WeakValidation)
	if err != nil {
		if strings.Contains(url, digestSep) {
			return &DigestedImage{nameWithDigest, fmt.Errorf("Expected valid digest reference, but found '%s', reason: %s", url, err)}
		}
		return nil // not a reference with digest
	}
	return &DigestedImage{nameWithDigest, nil}
}

