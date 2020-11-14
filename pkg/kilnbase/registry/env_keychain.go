// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	regauthn "github.com/google/go-containerregistry/pkg/authn"
	regname "github.com/google/go-containerregistry/pkg/name"
)

// Env vars with optional suffix:
//   export KILNBASE_REGISTRY_HOSTNAME_0=...
//   export KILNBASE_REGISTRY_USERNAME_0=...
//   export KILNBASE_REGISTRY_PASSWORD_0=...

type EnvKeychain struct {
	globalPrefix string

	infos       []envKeychainInfo
	collectErr  error
	collected   bool
	collectLock sync.Mutex
}

var _ regauthn.Keychain = &EnvKeychain{}

func NewEnvKeychain(globalPrefix string) *EnvKeychain {
	return &EnvKeychain{globalPrefix: globalPrefix}
}

func (k *EnvKeychain) Resolve(target regauthn.Resource) (regauthn.Authenticator, error) {
	infos, err := k.collect()
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		if strings.EqualFold(info.Hostname, target.RegistryStr()) {
			return regauthn.FromConfig(regauthn.AuthConfig{
				Username:      info.Username,
				Password:      info.Password,
				IdentityToken: info.IdentityToken,
				RegistryToken: info.RegistryToken,
			}), nil
		}
	}
	return regauthn.Anonymous, nil
}

type envKeychainInfo struct {
	Hostname      string
	Username      string
	Password      string
	IdentityToken string
	RegistryToken string
}

func (k *EnvKeychain) collect() ([]envKeychainInfo, error) {
	k.collectLock.Lock()
	defer k.collectLock.Unlock()

	if k.collected || k.collectErr != nil {
		return append([]envKeychainInfo{}, k.infos...), k.collectErr
	}

	const (
		sep = "_"
	)

	funcsMap := map[string]func(*envKeychainInfo, string) error{
		"HOSTNAME": func(info *envKeychainInfo, val string) error {
			registry, err := regname.NewRegistry(val, regname.StrictValidation)
			if err != nil {
				return errors.New("Parsing registry hostname: invalid value; " +
					"use a hostname such as gcr.io or index.docker.io " +
					"without URL credentials, query or fragment")
			}
			info.Hostname = registry.RegistryStr()
			return nil
		},
		"USERNAME": func(info *envKeychainInfo, val string) error {
			info.Username = val
			return nil
		},
		"PASSWORD": func(info *envKeychainInfo, val string) error {
			info.Password = val
			return nil
		},
		"IDENTITY_TOKEN": func(info *envKeychainInfo, val string) error {
			info.IdentityToken = val
			return nil
		},
