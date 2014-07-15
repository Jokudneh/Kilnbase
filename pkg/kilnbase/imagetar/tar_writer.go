// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package imagetar

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/imagedesc"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/util"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
)

type Logger interface {
	WriteStr(str string, args ...interface{}) error
}

type TarWriterOpts struct {
	Concurrency int
}

type TarWriter struct {
	ids       *imagedesc.ImageRefDescriptors
	dstOpener func() (io.WriteCloser, error)

	dst           io.WriteCloser
	tf            *tar.Writer
	layersToWrite []imagedesc.ImageLayerDescriptor

	opts   TarWriterOpts
	logger Logger
}

func NewTarWriter(ids *imagedesc.ImageRefDescriptors, dstOpener func() (io.WriteCloser, error),
	opts TarWriterOpts, logger Logger) *TarWriter {
	return &TarWriter{ids: ids, dstOpener: dstOpener, opts: opts, logger: logger}
}

func (w *TarWriter) Write() error {
	var err error

	w.dst, err = w.dstOpener()
	if err != nil {
		return err
	}

	defer w.dst.Close()

	w.tf = tar.NewWriter(w.dst)
	defer w.tf.Close()

	idsBytes, err := w.ids.AsBytes()
	if err != nil {
		return err
	}

	err = w.writeTarEntry(w.tf, "manifest.json", bytes.NewReader(idsBytes), int64(len(idsBytes)))
	if err != nil {
		return err
	}

	for _, td := range w.ids.Descriptors() {
		switch {
		case td.Image != nil:
			err := w.writeImage(*td.Image)
			if err != nil {
				return err
			}

		case td.ImageIndex != nil:
			err := w.writeImageIndex(*td.ImageIndex)
			if err != nil {
				return err
			}

		default:
			panic("Unknown item")
		}
	}

	return w.writeLayers()
}

func (w *TarWriter) writeImageIndex(td imagedesc.ImageIndexDescriptor) error {
	for _, idx := range td.Indexes {
		err := w.writeImageIndex(idx)
