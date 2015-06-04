// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"
	"regexp"

	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
)

type Matcher struct {
	url string
}

func NewMatcher(url string) Matcher { return Matcher{url} }

func (m Matcher) Matches(ref ctlconf.ImageRef) bool {
