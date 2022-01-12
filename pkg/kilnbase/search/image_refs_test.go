// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package search_test

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlser "github.com/Jokudneh/Kilnbase/pkg/kilnbase/search"
)

func TestImageRefsMatches(t *testing.T) {
	type matcherExample struct {
		InputResource  interface{}
		OutputResource interface{}
		SearchRules    []ctlconf.SearchRule
		OutputImages   []string
	}

	exs := []matcherExample{
		// By key
		{
			InputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "gcr.io/repo:something",
				},
			},
			OutputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "found:gcr.io/repo:something",
				},
			},
			SearchRules: []ctlconf.SearchRule{{
				KeyMatcher: &ctlconf.SearchRuleKeyMatcher{
					Name: "key",
				},
			}},
			OutputImages: []string{"gcr.io/repo:something"},
		},
		// By exact image
		{
			InputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "gcr.io/repo:something",
				},
			},
			OutputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "found:gcr.io/repo:something",
				},
			},
			SearchRules: []ctlconf.SearchRule{{
				ValueMatcher: &ctlconf.SearchRuleValueMatcher{
					Image: "gcr.io/repo:something",
				},
			}},
			OutputImages: []string{"gcr.io/repo:something"},
		},
		// By image repo
		{
			InputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
			},
			OutputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "found:gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
			},
			SearchRules: []ctlconf.SearchRule{{
				ValueMatcher: &ctlconf.SearchRuleValueMatcher{
					ImageRepo: "gcr.io/repo",
				},
			}},
			OutputImages: []string{"gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		},
		// By key and image repo
		{
			InputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key":   "gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
					"other": "gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
			},
			OutputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key":   "found:gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
					"other": "gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
			},
			SearchRules: []ctlconf.SearchRule{{
				KeyMatcher: &ctlconf.SearchRuleKeyMatcher{
					Name: "key",
				},
				ValueMatcher: &ctlconf.SearchRuleValueMatcher{
					ImageRepo: "gcr.io/repo",
				},
			}},
			OutputImages: []string{"gcr.io/repo@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		},
		// Matchers matching same values
		{
			InputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "gcr.io/repo:something",
				},
			},
			OutputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "found:gcr.io/repo:something",
				},
			},
			SearchRules: []ctlconf.SearchRule{{
				KeyMatcher: &ctlconf.SearchRuleKeyMatcher{
					Name: "key",
				},
			}, {
				KeyMatcher: &ctlconf.SearchRuleKeyMatcher{
					Name: "key",
				},
			}},
			OutputImages: []string{"gcr.io/repo:something"},
		},
		// JSON extraction
		{
			InputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": `{"nestedimage":"nginx1"}`,
				},
				"nestedimage": "nginx2",
			},
			OutputResource: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": `{"nestedimage":"found:nginx1"}`,
				},
				"nestedimage": "nginx2",
			},
			SearchRules: []ctlconf.SearchRule{{
				KeyMatcher: &ctlconf.SearchRuleKeyMatcher{
					Name: "key",
				},
				UpdateStrategy: &ctlconf.SearchRuleUpdateStrategy{
					JSON: &ctlconf.SearchRuleUpdateStrategyJSON{
						SearchRules: []ctlconf.SearchRule{{
							KeyMatcher: &ctlconf.SearchRuleKeyMatcher{
								Name: "nestedimage",
							}},
						},
					},
				},
			}},
			OutputImages: []string{"nginx1"},
		},
		// YAML extraction
		{
			InputResource: map[string]interface{}{
