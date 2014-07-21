// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"
	"strings"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	regname "github.com/google/go-containerregistry/pkg/name"
)

const digestSep = "@"
