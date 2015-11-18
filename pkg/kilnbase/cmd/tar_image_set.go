// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagetar"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
)

type TarImageSet struct {
	imageSet    ImageSet
	concurrency int
	logger      *ctllog.PrefixWriter
}

func (o TarImageSet) Export(foundImages *UnprocessedImageURLs,
