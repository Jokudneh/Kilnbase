// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image_test

import (
	"testing"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	. "github.com/Jokudneh/Kilnbase/pkg/kilnbase/image"
	"github.com/stretchr/testify/assert"
)

func TestPreResolvedOrigins(t *testing.T) {
	preresolvedImage := "image-resolved-with-sha@sha256:9999"
	originsFromImagesLockConfig := []ctlconf.Origin{{Preresolved: &ctlconf.OriginPreresolved{URL: preresolvedImage}}}
