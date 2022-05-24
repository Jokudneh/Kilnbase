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
  newImage: index.docker.io/cloudfoundry/cloud-controller-ng@sha256:374f967edd7db4d7efc2f38cb849988aa36a8248dd240d56f49484b8159fd800
  preresolved: true
- image: cloudfoundry/cnb:0.0.94-bionic@sha256:5b03a853e636b78c44e475bbc514e2b7b140cc41cca8ab907e9753431ae8c0b0
  newImage: index.docker.io/cloudfoundry/cnb@sha256:5b03a853e636b78c44e475bbc514e2b7b140cc41cca8ab907e9753431ae8c0b0
  preresolved: true
- image: docker.io/istio/citadel:1.4.5
  newImage: index.docker.io/istio/citadel@sha256:420a331a528886aca47bed5b8c549c78d594e52d771f876f3137d3557207712f
  preresolved: true
- image: docker.io/istio/galley:1.4.5
  newImage: index.docker.io/istio/galley@sha256:26e744bdfd3db289d4cfc9be63e38e7c7a424a9f76d1224cbdbbe58374229b68
  preresolved: true
- image: docker.io/istio/mixer:1.4.5
  newImage: index.docker.io/istio/mixer@sha256:ff6f39732c31999911790b00b484a471b6fe87192223d3266b0b6e752a374287
  preresolved: true
- image: docker.io/istio/node-agent-k8s:1.4.5
  newImage: index.docker.io/istio/node-agent-k8s@sha256:7e17ab509777a54f3c0dfb4518692a9ca179d1e8c41df87dc81a734339b37152
  preresolved: true
- image: docker.io/istio/pilot:1.4.5
  newImage: index.docker.io/istio/pilot@sha256:2bca5900d6bf20d5f0bf2b6673bc4d2885bab8cca2a9060336a0024930665b59
  preresolved: true
- image: docker.io/istio/proxyv2:1.4.5
  newImage: index.docker.io/istio/proxyv2@sha256:fc09ea0f969147a4843a564c5b677fbf3a6f94b56627d00b313b4c30d5fef094
  preresolved: true
- image: docker.io/istio/sidecar_injector:1.4.5
  newImage: index.docker.io/istio/sidecar_injector@sha256:ba446f8cf98bafdad4514fd492432dd180243cbc55a0b9c6bebfe31cb169033d
  preresolved: true
- image: eirini/opi@sha256:2e0b84c5fcb1e6e5cdb07a70210f2e462aa52119f7a330660a7444a938deefbb
