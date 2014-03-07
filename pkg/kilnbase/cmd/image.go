// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"reflect"
	"strings"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	"sigs.k8s.io/yaml"
)

type Images []Image

type Image struct {
	URL        string
	Origins    []ctlconf.Origin // empty when deserialized
