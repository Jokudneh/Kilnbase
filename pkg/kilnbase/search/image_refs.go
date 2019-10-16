// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"

	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
	"sigs.k8s.io/yaml"
)

type ImageRefs struct {
	res         interface{}
	searchRules []ctlconf.SearchRule
}

type ImageRefsVisitorFunc func(string) (string, bool)

func NewImageRefs(res interface{}, searchRules []ctlconf.SearchRule) ImageRefs {
	return ImageRefs{res, searchRules}
}

func (refs ImageRefs) Visit(visitorFunc ImageRefsVisitorFunc) {
	visitorFunc.Apply(refs.res, refs.searchRules)
}

func (v ImageRefsVisitorFunc) Apply(res interface{}, searchRules []ctlconf.SearchRule) {
	tmpRefs := map[string]string{}
	tmpRefPrefix := v.randomPrefix()
	tmpRefIdx := 0

	insertTmpRefsFunc := func(val string) (string, bool) {
		newVal, updated := v(val)
		if !updated {
			return "", false
		}
		// Only use temporary references for values that are updated;
		// otherwise future search rules may find tmp refs
		tmpRefID := fmt.Sprintf(tmpRefPrefix+"%d__", tmpRefIdx)
		tmpRefs[tmpRefID] = newVal
		tmpRefIdx++
		return tmpRefID, true
	}

	// Use a single matcher that represents all rules instead
	// so that each leaf value (string) is found once
	// even if it matches multiple search rules
	NewFields(res, RulesMatcher{searchRules}).Visit(v.extractValueFunc(insertTmpRefsFunc))

	resolveTmpRefsFunc := func(val string) (string, bool) {
		if actualRef, found := tmpRefs[val]; found {
			delete(tmpRefs, val)
			return actualRef, true
		}
		return "", false // TODO panic?
	}

	NewFields(res, tmpRefMatcher{tmpRefPrefix}).Visit(v.extractValueFunc(resolveTmpRefsFunc))

	if len(tmpRefs) > 0 {
		panic("ImageRefs: Expected all tmp refs to be found")
	}
}

func (v ImageRefsVisitorFunc) extractValueFunc(visitorFunc ImageRefsVisitorFunc) FieldsVisitorFunc {
	return func(val interface{}, ext ctlconf.SearchRuleUpdateStrategy) (interface{}, bool) {
		switch {
		case ext.None != nil:
			return val, false
