//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestPkgUnpkgSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	// redis:5.0.4
	input := `
kind: Object
spec:
# references image index
- image: index.docker.io/library/redis@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
# references plain image
- image: gcr.io/cloud-builders/gcs-fetcher@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
`

	path := "/tmp/kilnbase-test-pkg-unpkg-successful"
	defer os.RemoveAll(path)

	out, _ := kilnbase.RunWithOpts([]string{"package", "-f", "-", "--output", path}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := ""

	if out != expectedOut {
		t.Fatalf("Expected package output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}

	out, _ = kilnbase.RunWithOpts([]string{"unpackage", "-f", "-", "--input", path, "--repository", env.WithRegistries("docker.io/*username*/kilnbase-test-pkg-unpkg")}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut = env.WithRegistries(`---
kind: Object
spec:
- image: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
- image: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
`)

	if out != expectedOut {
		t.Fatalf("Expected unpackage output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestPkgUnpkgLockSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
overrides:
# ignored because it's not preresolved
- image: gcs-fetcher
  newImage: gcr.io/cloud-builders/gcs-fetcher@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
- image: redis
  newImage: index.docker.io/library/redis@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
  preresolved: true
`

	path := "/tmp/kilnbase-test-pkg-unpkg-with-lock-successful"
	defer os.RemoveAll(path)

	relocatedLockPath := path + "-relocated"
	defer os.RemoveAll(relocatedLockPath)

	out, _ := kilnbase.RunWithOpts([]string{"package", "-f", "-", "--output", path}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := ""

	if out != expectedOut {
		t.Fatalf("Expected package output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}

	out, _ = kilnbase.RunWithOpts([]string{
		"unpackage", "-f", "-", "--input", path, "--repository",
		env.WithRegistries("docker.io/*username*/kilnbase-test-pkg-unpkg"),
		"--lock-output", relocatedLockPath,
	}, RunOpts{StdinReader: strings.NewReader(input)})

	lockOutBs, err := ioutil.ReadFile(relocatedLockPath)
	if err != nil {
		t.Fatalf("Expected to find relocated lock file")
	}

	kilnbaseVersionOutput, _ := kilnbase.RunWithOpts([]string{"version"}, RunOpts{})
	kilnbaseVersion := strings.SplitAfter(kilnbaseVersionOutput, " ")[2]

	expectedLockOut := strings.ReplaceAll(env.WithRegistries(`apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
minimumRequiredVersion: __ver__
overrides:
- image: redis
  newImage: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
  preresolved: true
- image: index.docker.io/library/redis@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
  newImage: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
  preresolved: true
`), "__ver__", strings.TrimSuffix(kilnbaseVersion, "\n"))

	if string(lockOutBs) != expectedLockOut {
		t.Fatalf("Expected unpackage lock output >>>%s<<< to match >>>%s<<<", lockOutBs, expectedLockOut)
	}
}

func TestPkgUnpkgSuccessfulWithForeignLayers(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	// Mongo has 2 foreign layers
	input := `
kind: Object
spec:
- image: index.docker.io/library/mongo@sha256:633ec3ae6db954a65a1abadb482bae73375d0098005cb36a3851b32cd891b22e
`

	path := "/tmp/kilnbase-test-pkg-unpkg-successful-foreign-layers"
	defer os.RemoveAll(path)

	out, _ := kilnbase.RunWithOpts([]string{"package", "-f", "-", "--output", path}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := ""

	if out != expectedOut {
		t.Fatalf("Expected package output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}

	out, _ = kilnbase.RunWithOpts([]string{"unpackage", "-f", "-", "--input", path, "--repository", env.WithRegistries("docker.io/*username*/kilnbase-test-pkg-unpkg")}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut = env.WithRegistries(`---
kind: Object
spec:
- image: index.docker.io/*username*/kilnbase-test-pkg-unpkg@sha256:633ec3ae6db954a65a1abadb482bae73375d0098005cb36a3851b32cd891b22e
`)

	if out != expectedOut {
		t.Fatalf("Expected unpackage output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestPkgUnpkgSuccessfulWithManyImages(t *testing.T) {
	env := BuildEnv(t)

	if env.SkipStressTests {
		fmt.Printf("This is a stress test; skipping.")
		return
	}

	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
minimumRequiredVersion: 0.20.0
overrides:
- image: cfidentity/uaa@sha256:9f1e7e399c96309935145624d1824b2c2bf93656fd9c4dcf1c593b55f98aa6a8
  newImage: index.docker.io/cfidentity/uaa@sha256:9f1e7e399c96309935145624d1824b2c2bf93656fd9c4dcf1c593b55f98aa6a8
  preresolved: true
- image: cloudfoundry/capi:nginx@sha256:51e4e48c457d5cb922cf0f569e145054e557e214afa78fb2b312a39bb2f938b6
  newImage: index.docker.io/cloudfoundry/capi@sha256:51e4e48c457d5cb922cf0f569e145054e557e214afa78fb2b312a39bb2f938b6
  preresolved: true
- image: cloudfoundry/cloud-controller-ng:1ebab1cbb5270a3d51c0a098a37cd9e8ca0f721d@sha256:374f967edd7db4d7efc2f38cb849988aa36a8248dd240d56f49484b8159fd800
