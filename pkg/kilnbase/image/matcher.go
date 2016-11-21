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
	switch {
	case len(ref.Image) > 0:
		// Registry hostnames are case-insensitive (RFC 1035/1123); ref.Image
		// comes from scanning the target manifest while m.url comes from the
		// user's separately-authored kilnbase.yaml config, so a hostname-case
		// mismatch between the two sources must not prevent a match.
		lhs := ctlbdk.LowerCaseRegistry(ref.Image)
		rhs := ctlbdk.LowerCaseRegistry(m.url)
		return lhs == rhs

	case len(ref.ImageRepo) > 0:
		repo, _ := URLRepo(m.url)
		lhs := ctlbdk.LowerCaseRegistry(ref.ImageRepo)
		rhs := ctlbdk.LowerCaseRegistry(repo)
