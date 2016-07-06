// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"fmt"
	"os"

	"kilnbase.dev/imgpkg/pkg/imgpkg/lockconfig"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	"github.com/Jokudneh/Kilnbase/pkg/kilnbase/version"
	versions "kilnbase.dev/vendir/pkg/vendir/versions/v1alpha1"
	semver "github.com/hashicorp/go-version"
	"sigs.k8s.io/yaml"
)

const (
	configAPIVersion      = "kilnbase.k14s.io/v1alpha1"
	configKind            = "Config"
	sourcesKind           = "Sources"           // specify list of sources for building images
	imageOverridesKind    = "ImageOverrides"    // specify alternative image urls
	imageDestinationsKind = "ImageDestinations" // specify image push destinations
	imageKeysKind         = "ImageKeys"
)

type Kind struct {
	APIVersion, Kind string
}

var (
	configKinds = []Kind{
		{configAPIVersion, configKind},
		{configAPIVersion, configKind},
		{configAPIVersion, sourcesKind},
		{configAPIVersion, imageOverridesKind},
		{configAPIVersion, imageDestinationsKind},
		{configAPIVersion, imageKeysKind},
	}
)

type Config struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind,omitempty"`

	MinimumRequiredVersion string `json:"minimumRequiredVersion,omitempty"`

	Sources      []Source           `json:"sources,omitempty"`
	Overrides    []ImageOverride    `json:"overrides,omitempty"`
	Destinations []ImageDestination `json:"destinations,omitempty"`
	Keys         []string           `json:"keys,omitempty"`
	SearchRules  []SearchRule       `json:"searchRules,omitempty"`
}

type Source struct {
	ImageRef
	Path string

	Docker          *SourceDockerOpts
	Pack            *SourcePackOpts
	KubectlBuildkit *SourceKubectlBuildkitOpts
	Ko              *SourceKoOpts
	Bazel           *SourceBazelOpts
	Maven           *SourceJibOpts
}

type ImageOverride struct {
	ImageRef
	NewImage          string                     `json:"newImage"`
	Preresolved       bool                       `json:"preresolved,omitempty"`
	TagSelection      *versions.VersionSelection `json:"tagSelection,omitempty"`
	PlatformSelection *PlatformSelection         `json:"platformSelection,omitempty"`
	ImageOrigins      []Origin                   `json:"origins,omitempty"`
}

// PlatformSelection
// https://github.com/opencontainers/image-spec/blob/main/image-index.md#image-index-property-descriptions
type PlatformSelection struct {
	Architecture string   `json:"architecture,omitempty"`
	OS           string   `json:"os,omitempty"`
	OSVersion    string   `json:"os.version,omitempty"`  // note dot
	OSFeatures   []string `json:"os.features,omitempty"` // note dot
	Variant      string   `json:"variant,omitempty"`
	Features     []string `json:"features,omitempty"`
}

type ImageDestination struct {
	ImageRef
	NewImage string   `json:"newImage"`
	Tags     []string `json:"tags"`
}

type SearchRule struct {
	KeyMatcher     *SearchRuleKeyMatcher     `json:"keyMatcher,omitempty"`
	ValueMatcher   *SearchRuleValueMatcher   `json:"valueMatcher,omitempty"`
	UpdateStrategy *SearchRuleUpdateStrategy `json:"updateStrategy,omitempty"`
	// TODO ResourceMatchers (see kapp's matchers)
}

type SearchRuleKeyMatcher struct {
	Name string      `json:"name,omitempty"`
	Path ctlres.Path `json:"path,omitempty"`
	// TODO JSONPath string
}

type SearchRuleValueMatcher struct {
	Image     string `json:"image,omitempty"`
	ImageRepo string `json:"imageRepo,omitempty"`
	// TODO Regexp    string `json:"regexp,omitempty"`
}

type SearchRuleUpdateStrategy struct {
