// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package search

import (
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
)

type RulesMatcher struct {
	rules []ctlconf.SearchRule
}

func NewRulesMatcher(rules []ctlconf.SearchRule) RulesMatcher {
	return RulesMatcher{rules}
}

var _ Matcher = RuleMatcher{}

func (m RulesMatcher) Matches(keyPath ctlres.Path, value interface{}) (bool, ctlconf.SearchRuleUpdateStrategy) {
	for _, rule := range m.rules {
		matches, extraction := (RuleMatcher{rule}).Matches(keyPath, value)
		if matches {
			return true, extraction
		}
