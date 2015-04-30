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

