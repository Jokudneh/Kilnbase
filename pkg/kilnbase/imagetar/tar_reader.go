// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagetar

import (
	"io"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagedesc"
)

type TarReader struct {
	path string
}

func NewTarReader(path string) TarReader {
	return TarReader{path}
}

func (r TarReader) Read() ([]imagedesc.ImageOrIndex, error) {
	file := tarFile{r.path}

	manifestFile, err := file.Chunk("manifest.json").Open()
	if err != nil {
