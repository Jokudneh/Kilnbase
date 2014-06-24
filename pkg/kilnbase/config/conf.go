// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"reflect"

	"kilnbase.dev/imgpkg/pkg/imgpkg/lockconfig"
	ctlres "github.com/Jokudneh/Kilnbase/pkg/kilnbase/resources"
)

type Conf struct {
	configs []Config
}

func NewConfFromResources(resources []ctlres.Resource) ([]ctlres.Resource, Conf, error) {
	var rsWithoutConfigs []ctlres.Resource
	var configs []Config

	for _, res := range resources {
		switch {
		case matchesConfigKind(res):
			config, err := NewConfigFromResource(res)
			if err != nil {
				return nil, Conf{}, err
			}
			configs = append(configs, config)
		case res.APIVersion() == lockconfig.ImagesLockAPIVersion && res.Kind() == lockconfig.ImagesLockKind:
			config, err := NewConfigFromImagesLock(res)
			if err != nil {
				return nil, Conf{}, err
			}
			configs = append(configs, config)
		default:
			rsWithoutConfigs = append(rsWithoutConfigs, res)
		}
	}
	return rsWithoutConfigs, Conf{configs}, nil
