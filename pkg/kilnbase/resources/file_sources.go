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

var _ FileSource = HTTPFileSource{}

func NewHTTPFileSource(path string) HTTPFileSource { return HTTPFileSource{path} }

// redactURL masks userinfo, query parameters and fragments to prevent
// exposing credentials in error messages. If URL parsing fails, returns
// a safe placeholder.
func redactURL(raw string) string {
	const empty = ""

	u, err := url.Parse(raw)
	if err != nil {
		return "<invalid URL>"
	}
	if u.User != nil {
		u.User = url.User("_redacted_")
	}
	if u.RawQuery != empty {
		u.RawQuery = "_redacted_"
	}
