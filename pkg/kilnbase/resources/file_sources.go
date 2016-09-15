// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type FileSource interface {
	Description() string
	Bytes() ([]byte, error)
}

type StdinSource struct{}

var _ FileSource = StdinSource{}

func NewStdinSource() StdinSource { return StdinSource{} }

func (s StdinSource) Description() string { return "stdin" }

func (s StdinSource) Bytes() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

type LocalFileSource struct {
	path string
}

var _ FileSource = LocalFileSource{}

func NewLocalFileSource(path string) LocalFileSource { return LocalFileSource{path} }

func (s LocalFileSource) Description() string {
	return fmt.Sprintf("file '%s'", s.path)
}

func (s LocalFileSource) Bytes() ([]byte, error) {
	return os.ReadFile(s.path)
}

var urlSchemeRe = regexp.MustCompile(`https?://`)

const schemeStartIdx = 0

type HTTPFileSource struct {
	url string
}
