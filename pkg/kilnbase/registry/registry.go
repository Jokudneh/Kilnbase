// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	regauthn "github.com/google/go-containerregistry/pkg/authn"
	regname "github.com/google/go-containerregistry/pkg/name"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
	regremote "github.com/google/go-containerregistry/pkg/v1/remote"
)

type Opts struct {
	CACertPaths   []string
	VerifyCerts   bool
	Insecure      bool
	EnvAuthPrefix string
}

type Registry struct {
	opts    []regremote.Option
	refOpts []regname.Option
}

func NewRegistry(opts Opts) (Registry, error) {
	keychain := regauthn.NewMultiKeychain(NewEnvKeychain(opts.EnvAuthPrefix), regauthn.DefaultKeychain)
	transport, err := newHTTPTransport(opts)
	if err != nil {
		return Registry{}, err
	}

	var refOpts []regname.Option
	if opts.Insecure {
		refOpts = append(refOpts, regname.Insecure)
	}

	return Registry{
		opts: []regremote.Option{
			regremote.WithTransport(transport),
			regremote.WithAuthFromKeychain(keychain),
		},
		refOpts: refOpts,
	}, nil
}

func (i Registry) Generic(ref regname.Reference) (regv1.Descriptor, error) {
	ref, err := regname.ParseReference(ref.String(), i.refOpts...)
	if err != nil {
		return regv1.Descriptor{}, err
	}

	desc, err := regremote.Head(ref, i.opts...)
	if err != nil {
		getDesc, err := regremote.Get(ref, i.opts...)
		if err != nil {
			return regv1.Descriptor{}, err
		}
		return getDesc.Descriptor, err
	}

	return *desc, nil
}

func (i Registry) Image(ref regname.Reference) (regv1.Image, error) {
	ref, err := regname.ParseReference(ref.String(), i.refOpts...)
	if err != nil {
		return nil, err
	}

	return regremote.Image(ref, i.opts...)
}

func (i Registry) WriteImage(ref regname.Reference, img regv1.Image) error {
	ref, err := regname.ParseReference(ref.String(), i.refOpts...)
	if err != nil {
		return err
	}

	err = i.retry(func() error {
		return regremote.Write(ref, img, i.opts...)
	})
	if err != nil {
		return fmt.Errorf("Writing image: %s", err)
	}

	return nil
}

func (i Registry) Index(ref regname.Reference) (regv1.ImageIndex, error) {
	ref, err := regname.ParseReference(ref.String(), i.refOpts...)
	if err != nil {
		return nil, err
	}

	return regremote.Index(ref, i.opts...)
}

func (i Registry) WriteIndex(ref regname.Reference, idx regv1.ImageIndex) error {
	ref, err := regname.ParseReference(ref.String(), i.refOpts...)
	if err != nil {
		return err
	}

	err = i.retry(func() error {
		return regremote.WriteIndex(ref, idx, i.opts...)
	})
	if err != nil {
		return fmt.Errorf("Writing image index: %s", err)
	}

	return nil
}

func (i Registry) WriteTag(dstRef regname.Tag, srcRef regname.Digest) error {
	dstRef, err := regname.NewTag(dstRef.String(), i.refOpts...)
	if err != nil {
		return err
	}

	srcRef, err = regname.NewDigest(srcRef.String(), i.refOpts...)
	if err != nil {
		return err
	}

	err = i.retry(func() error {
		desc, err := regremote.Get(srcRef, i.opts...)
		if err != nil {
			return err
		}

		return regremote.Tag(dstRef, desc, i.opts...)
	})
	if err != nil {
		return fmt.Errorf("Writing image tag: %s", err)
	}

	return nil
}

