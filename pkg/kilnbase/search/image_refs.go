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
