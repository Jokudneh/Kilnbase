// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"fmt"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlreg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/registry"
	"kilnbase.dev/vendir/pkg/vendir/versions"
	"kilnbase.dev/vendir/pkg/vendir/versions/v1alpha1"
	regname "github.com/google/go-containerregistry/pkg/name"
)

// TagSelectedImage represents an image that will be resolved into url+digest
