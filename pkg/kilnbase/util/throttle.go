// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
)

type Throttle struct {
	ch chan struct{}
}

