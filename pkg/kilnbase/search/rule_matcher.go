// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"reflect"

	ctlbdk "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder/docker"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlimg "github.com/Jokudneh/Kilnbase/pkg/kilnbase/image"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
)

type Matcher interface {
	Matches(keyPath ctlres.Path, value interface{}) (bool, ctlconf.SearchRuleUpdateStrategy)
}

type RuleMatcher struct {
	rule ctlconf.SearchRule
}

var _ Matcher = RuleMatcher{}

func (m RuleMatcher) Matches(keyPath ctlres.Path, value interface{}) (bool, ctlconf.SearchRuleUpdateStrategy) {
	var keyMatched, valueMatched bool

	if m.rule.KeyMatcher != nil {
		switch {
		case len(m.rule.KeyMatcher.Name) > 0:
			name := m.rule.KeyMatcher.Name
			keyMatched = keyPath.HasMatchingSuffix(ctlres.Path{{MapKey: &name}})
