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

