// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	tagBuilderTagCleanRegexp = regexp.MustCompile("[^a-zA-Z0-9\\-]+")
)

type TagBuilder struct{}

func (d TagBuilder) CheckTagLen128(tag string) string {
	// "A tag ... may contain a maximum of 128 characters."
	// (https://docs.docker.com/engine/reference/commandline/tag/)
	return d.CheckLen(tag, 128)
}

func (d TagBuilder) CheckLen(str string, num int) string {
	if len(str) > num {
		panic(fmt.Sprintf("Expected string '%s' len to be less than %d", str, num))
	}
	return str
}

