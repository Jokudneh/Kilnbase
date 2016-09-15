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
