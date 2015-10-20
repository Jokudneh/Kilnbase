// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
)

type Throttle struct {
	ch chan struct{}
}

func NewThrottle(max int) Throttle {
	if max < 1 {
		panic(fmt.Sprintf("Expected maximum throttle to be >= 1, but was %d", max))
	}
	ch := make(chan struct{}, max)
