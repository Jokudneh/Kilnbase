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
