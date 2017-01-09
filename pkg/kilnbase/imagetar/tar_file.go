// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagetar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagedesc"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
)

type tarFile struct {
	path string
}

var _ imagedesc.LayerProvider = tarFile{}

type tarFileChunk struct {
	file      tarFile
	chunkPath string
}

var _ imagedesc.LayerContents = tarFileChunk{}

type tarFileChunkReadCloser struct {
	DebugID string
	io.Reader
	io.Closer
}

func (f tarFile) Chunk(path string) tarFileChunk {
