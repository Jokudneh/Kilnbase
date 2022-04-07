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
		if err != nil {
			return err
		}
	}

	for _, img := range td.Images {
		err := w.writeImage(img)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *TarWriter) writeImage(td imagedesc.ImageDescriptor) error {
	for _, imgLayer := range td.Layers {
		// TODO anything else we can do to deal with this?
		// Do not include foreign layers since we cannot
		// import without changing image digest.
		if imgLayer.IsDistributable() {
			w.layersToWrite = append(w.layersToWrite, imgLayer)
		}
	}
	return nil
}

type writtenLayer struct {
	Name   string
	Offset int64
	Layer  imagedesc.ImageLayerDescriptor
}

func (w *TarWriter) writeLayers() error {
	// Sort layers by digest to have deterministic archive
	sort.Slice(w.layersToWrite, func(i, j int) bool {
		return w.layersToWrite[i].Digest < w.layersToWrite[j].Digest
	})

	seekableDst, isSeekable := w.dst.(*os.File)
	isInflattable := (w.opts.Concurrency > 1) && isSeekable
	writtenLayers := map[string]writtenLayer{}

	// Inflate tar file so that multiple writes can happen in parallel
	for _, imgLayer := range w.layersToWrite {
		digest, err := regv1.NewHash(imgLayer.Digest)
		if err != nil {
			return err
		}

		name := digest.Algorithm + "-" + digest.Hex + ".tar.gz"

		// Dedup layers
		if _, found := writtenLayers[name]; found {
			continue
		}

		err = w.tf.Flush()
		if err != nil {
			return err
		}

		var stream io.Reader
		var currPos int64

		if isSeekable {
			currPos, err = seekableDst.Seek(0, 1)
			if err != nil {
				return fmt.Errorf("Finding current position: %s", err)
			}
		}

		if isInflattable {
			stream = nil
		} else {
			foundLayer, err := w.ids.FindLayer(imgLayer)
			if err != nil {
				return err
			}

			stream, err = foundLayer.Open()
			if err != nil {
				return err
			}
		}

		err = w.writeTarEntry(w.tf, name, stream, imgLayer.Size)
		if err != nil {
			return fmt.Errorf("Writing tar entry: %s", err)
		}

		writtenLayers[name] = writtenLayer{
			Name:   name,
			Layer:  imgLayer,
			Offset: currPos,
		}
	}

	err := w.tf.Flush()
	if err != nil {
		return err
	}

	if isInflattable {
		return w.fillInLayers(writtenLayers)
	}

	return nil
}

func (w *TarWriter) fillInLayers(writtenLayers map[string]writtenLayer) error {
	var sortedWritterLayers []writtenLayer

	for _, writtenLayer := range writtenLayers {
		sortedWritterLayers = append(sortedWritterLayers, writtenLayer)
	}

	// Prefer larger sizes first
	sort.Slice(sortedWritterLayers, func(i, j int) bool {
		return sortedWritterLayers[i].Layer.Size >= sortedWritterLayers[j].Layer.Size
	})

	errCh := make(chan error, len(writtenLayers))
	writeThrottle := util.NewThrottle(w.opts.Concurrency)

	// Fill in actual data
	for _, writtenLayer := range sortedWritterLayers {
		writtenLayer := writtenLayer // copy

		go func() {
			writeThrottle.Take()
			defer writeThrottle.Done()

			errCh <- w.retry(func() error {
				return w.fillInLayer(writtenLayer)
			})
		}()
	}

	for i := 0; i < len(writtenLayers); i++ {
		err := <-errCh
