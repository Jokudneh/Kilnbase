//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRelocateSuccessful(t *testing.T) {
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

	out, _ := kilnbase.RunWithOpts([]string{"relocate", "-f", "-", "--repository", env.WithRegistries("docker.io/*username*/kilnbase-test-relocate")}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := env.WithRegistries(`---
kind: Object
spec:
- image: index.docker.io/*username*/kilnbase-test-relocate@sha256:000339fb57e0ddf2d48d72f3341e47a8ca3b1beae9bdcb25a96323095b72a79b
- image: index.docker.io/*username*/kilnbase-test-relocate@sha256:055519529bf1ba12bf916fa42d6d3f68bdc581413621c269425bb0fee2467a93
`)

	if out != expectedOut {
		t.Fatalf("Expected unpackage output >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}
}

func TestRelocateSuccessfulWithManyImages(t *testing.T) {
	env := BuildEnv(t)

	if env.SkipStressTests {
		fmt.Printf("This is a stress test; skipping.")
		return
	}

	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: index.docker.io/cfidentity/uaa@sha256:9f1e7e399c96309935145624d1824b2c2bf93656fd9c4dcf1c593b55f98aa6a8
- image: index.docker.io/cloudfoundry/capi@sha256:51e4e48c457d5cb922cf0f569e145054e557e214afa78fb2b312a39bb2f938b6
- image: index.docker.io/cloudfoundry/cloud-controller-ng@sha256:374f967edd7db4d7efc2f38cb849988aa36a8248dd240d56f49484b8159fd800
- image: index.docker.io/cloudfoundry/cnb@sha256:5b03a853e636b78c44e475bbc514e2b7b140cc41cca8ab907e9753431ae8c0b0
- image: index.docker.io/istio/citadel@sha256:420a331a528886aca47bed5b8c549c78d594e52d771f876f3137d3557207712f
- image: index.docker.io/istio/galley@sha256:26e744bdfd3db289d4cfc9be63e38e7c7a424a9f76d1224cbdbbe58374229b68
- image: index.docker.io/istio/mixer@sha256:ff6f39732c31999911790b00b484a471b6fe87192223d3266b0b6e752a374287
- image: index.docker.io/istio/node-agent-k8s@sha256:7e17ab509777a54f3c0dfb4518692a9ca179d1e8c41df87dc81a734339b37152
- image: index.docker.io/istio/pilot@sha256:2bca5900d6bf20d5f0bf2b6673bc4d2885bab8cca2a9060336a0024930665b59
- image: index.docker.io/istio/proxyv2@sha256:fc09ea0f969147a4843a564c5b677fbf3a6f94b56627d00b313b4c30d5fef094
