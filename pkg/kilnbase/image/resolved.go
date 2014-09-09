// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	regname "github.com/google/go-containerregistry/pkg/name"
)

// ResolvedImage represents an image that will be resolved into url+digest
