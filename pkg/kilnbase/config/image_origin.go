// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"sigs.k8s.io/yaml"
)

type Origin struct {
	Git              *OriginGit              `json:"git,omitempty"`
	Local            *OriginLocal            `json:"local,omitempty"`
	Resolved         *OriginResolved         `json:"resolved,omitempty"`
	Tagged           *OriginTagged           `json:"tagged,omitempty"`
	Preresolved      *OriginPreresolved      `json:"preresolved,omitempty"`
	PlatformSelected *OriginPlatformSelected `json:"platformSelected,omitempty"`
}

