// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
)

type ErrImage struct {
	err error
}

var _ Image = ErrImage{}
