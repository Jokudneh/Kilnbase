// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	fileResourcesAllowedExts = []string{".json", ".yaml", ".yml"} // matches kubectl
)

type FileResource struct {
	fileSrc FileSource
}

func NewFileResources(file string) ([]FileResource, error) {
	var fileRs []FileResource

	switch {
	case file == "-":
		fileRs = append(fileRs, FileResource{NewStdinSource()})

	case strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://"):
		fileRs = append(fileRs, FileResource{NewHTTPFileSource(file)})

	default:
		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("Unable to stat file: %s", err)
		}

		if fileInfo.IsDir() {
			var paths []string

			err := filepath.Walk(file, func(path string, fi os.FileInfo, err error) error {
				if err != nil || fi.IsDir() {
