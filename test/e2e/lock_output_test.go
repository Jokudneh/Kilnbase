//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const failedReading = "Failed while reading %s"

var (
	imgLockWithResolvedOrigins = `---
apiVersion: imgpkg.kilnbase.dev/v1alpha1
images:
- annotations:
    kilnbase.kilnbase.dev/id: nginx:1.14.2
    kilnbase.kilnbase.dev/origins: |
      - resolved:
          tag: 1.14.2
          url: nginx:1.14.2
  image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- annotations:
    kilnbase.kilnbase.dev/id: sample-app
    kilnbase.kilnbase.dev/origins: |
      - resolved:
          tag: 1.15.1
          url: nginx:1.15.1
  image: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
kind: ImagesLock
`
	imgLockWithBuiltOrigins = `---
apiVersion: imgpkg.kilnbase.dev/v1alpha1
images:
- annotations:
    kilnbase.kilnbase.dev/id: nginx:1.14.2
    kilnbase.kilnbase.dev/origins: |
      - local:
          path: path/to/source
      - git:
          dirty: true
          remoteURL: git@github.com:jokudneh/kilnbase-kilnbase.git
          sha: f7988fb6c02e0ce69257d9bd9cf37ae20a60f1d
  image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- annotations:
    kilnbase.kilnbase.dev/id: sample-app
    kilnbase.kilnbase.dev/origins: |
      - local:
          path: path/to/source
      - git:
          dirty: true
          remoteURL: git@github.com:jokudneh/kilnbase-kilnbase.git
          sha: 4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1
  image: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
kind: ImagesLock
`
	imgLockWithBuiltAndPreresolvedOrigins = `---
apiVersion: imgpkg.kilnbase.dev/v1alpha1
images:
- annotations:
    kilnbase.kilnbase.dev/id: nginx:1.14.2
    kilnbase.kilnbase.dev/origins: |
      - local:
          path: path/to/source
      - git:
          dirty: true
          remoteURL: git@github.com:jokudneh/kilnbase-kilnbase.git
          sha: f7988fb6c02e0ce69257d9bd9cf37ae20a60f1d
      - preresolved:
          url: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
  image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- annotations:
    kilnbase.kilnbase.dev/id: sample-app
    kilnbase.kilnbase.dev/origins: |
      - local:
          path: path/to/source
      - git:
          dirty: true
          remoteURL: git@github.com:jokudneh/kilnbase-kilnbase.git
          sha: 4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1
      - preresolved:
          url: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
  image: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
kind: ImagesLock
`
	imgLock = `---
apiVersion: imgpkg.kilnbase.dev/v1alpha1
images:
- annotations:
    kilnbase.kilnbase.dev/id: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
  image: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
- annotations:
    kilnbase.kilnbase.dev/id: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
  image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
kind: ImagesLock
`
)

func TestLockOutputSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
images:
- image: nginx:1.14.2
- image: sample-app
- sidecarImage: sample-app
- - sample-app
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageOverrides
overrides:
- image: sample-app
  newImage: nginx:1.15.1
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageKeys
keys:
- sidecarImage
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
searchRules:
- keyMatcher:
    path: [images, {allIndexes: true}, {index: 0}]
`

	path := "/tmp/kilnbase-test-lock-output-successful"
	defer os.RemoveAll(path)

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false", "--lock-output=" + path}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	kilnbaseVersionOutput, _ := kilnbase.RunWithOpts([]string{"version"}, RunOpts{})
	kilnbaseVersion := strings.SplitAfter(kilnbaseVersionOutput, " ")[2]

	expectedOut := `---
images:
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
- sidecarImage: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
- - index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
`
	if out != expectedOut {
		t.Fatalf("Expected >>>%s<<< to match >>>%s<<<", out, expectedOut)
	}

	expectedFileContents := strings.ReplaceAll(`apiVersion: kilnbase.k14s.io/v1alpha1
kind: Config
minimumRequiredVersion: __ver__
overrides:
- image: nginx:1.14.2
  newImage: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
  preresolved: true
- image: sample-app
  newImage: index.docker.io/library/nginx@sha256:4a5573037f358b6cdfa2f3e8a9c33a5cf11bcd1675ca72ca76fbe5bd77d0d682
  preresolved: true
searchRules:
- keyMatcher:
    name: sidecarImage
- keyMatcher:
    path:
    - images
    - allIndexes: true
    - index: 0
`, "__ver__", strings.TrimSuffix(kilnbaseVersion, "\n"))

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf(failedReading, path)
	}
