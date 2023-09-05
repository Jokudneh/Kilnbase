//go:build e2e

// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveSuccessful(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: nginx:1.14.2
- image: library/nginx:1.14.2
- image: docker.io/library/nginx:1.14.2
- image: index.docker.io/library/nginx:1.14.2
- image: nginx@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
`

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
spec:
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
`

	require.YAMLEq(t, expectedOut, out)
}

func TestResolveSuccessfulWithAnnotations(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	// The repetition in this input is so it can test for:
	// 1) resolving
	// 2) filtering annotations with null origins (which happens with digests, which don't get resolved)
	// 3) de-duplicating annotations
	input := `
kind: Object
spec:
- image: nginx:1.14.2
- image: nginx:1.14.2
- image: library/nginx:1.14.2
- image: docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
`

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
metadata:
  annotations:
    kilnbase.k14s.io/images: |
      - origins:
        - resolved:
            tag: 1.14.2
            url: nginx:1.14.2
        url: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
      - origins:
        - resolved:
            tag: 1.14.2
            url: library/nginx:1.14.2
        url: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
spec:
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
`

	require.YAMLEq(t, expectedOut, out)
}

func TestSortAnnotations(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: foo/img1:bbb
- image: foo/img1:aaa
- image: foo/img1:ccc
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageOverrides
overrides:
- image: foo/img1:bbb
  newImage: bbb
  preresolved: true
- image: foo/img1:aaa
  newImage: aaa
  preresolved: true
- image: foo/img1:ccc
  newImage: ccc
  preresolved: true
`

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
metadata:
  annotations:
    kilnbase.k14s.io/images: |
      - origins:
        - preresolved:
            url: aaa
        url: aaa
      - origins:
        - preresolved:
            url: bbb
        url: bbb
      - origins:
        - preresolved:
            url: ccc
        url: ccc
spec:
- image: bbb
- image: aaa
- image: ccc
`

	require.YAMLEq(t, expectedOut, out)
}

func TestResolveInvalidDigest(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: nginx@sha256:digest
`

	var stderr bytes.Buffer
	out, err := kilnbase.RunWithOpts([]string{"--json", "-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader:  strings.NewReader(input),
		StderrWriter: &stderr,
		AllowError:   true,
	})

	require.EqualError(t, err, "Execution error: stdout: '' stderr: '' error: 'exit status 1'")
	require.Empty(t, out)
	require.Equal(t, "kilnbase: Error: \n- Resolving image 'nginx@sha256:digest': Expected valid digest reference, but found 'nginx@sha256:digest', reason: invalid checksum digest length\n", stderr.String())
}

func TestResolveUnknownImage(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: unknown
`

	_, err := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
		AllowError:  true,
	})

	expectedErr := "- Resolving image 'unknown': GET https://index.docker.io/v2/library/unknown/manifests/latest: UNAUTHORIZED: authentication required;"
	require.Contains(t, err.Error(), expectedErr)
}

func TestResolveWithOverride(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: unknown
- image: final
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageOverrides
overrides:
- image: unknown
  newImage: docker.io/library/nginx:1.14.2
- image: final
  newImage: docker.io/library/nginx:1.14.2
  preresolved: true
`

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
spec:
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: docker.io/library/nginx:1.14.2
`

	require.YAMLEq(t, expectedOut, out)
}

func TestResolveWithImageMap(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: img1
- image: img2:1.14.2
- image: img3
`

	imageMapData := `{
  "img1": "docker.io/foo/img1@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d",
  "img2:1.14.2": "docker.io/foo/img2:1.14.2",
  "img3": "img3"
}
`

	file, err := ioutil.TempFile("", "kilnbase-test-resolve-with-image-map")
	if err != nil {
		t.Fatalf("temp file err: %s", err)
	}

	file.Close()
	defer os.RemoveAll(file.Name())

	err = ioutil.WriteFile(file.Name(), []byte(imageMapData), os.ModePerm)
	if err != nil {
		t.Fatalf("write image map err: %s", err)
	}

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false", "--image-map-file", file.Name()}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
spec:
- image: docker.io/foo/img1@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: docker.io/foo/img2:1.14.2
- image: img3
`

	require.YAMLEq(t, expectedOut, out)
}

func TestResolveWithImageKeys(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: nginx:1.14.2
- customImage: nginx:1.14.2
- subPath:
    anotherCustomImage: nginx:1.14.2
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageKeys
keys:
- customImage
- anotherCustomImage
`

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
spec:
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- customImage: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- subPath:
    anotherCustomImage: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
`

	require.YAMLEq(t, expectedOut, out)
}

func TestResolveWithOverrideMatchingImageRepo(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: docker.io/library/not-nginx
- image: docker.io/library/not-nginx:1.14.2
- image: docker.io/library/not-nginx:1.20@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
---
apiVersion: kilnbase.k14s.io/v1alpha1
kind: ImageOverrides
overrides:
- imageRepo: docker.io/library/not-nginx
  newImage: docker.io/library/nginx:1.14.2
`

	out, _ := kilnbase.RunWithOpts([]string{"-f", "-", "--images-annotation=false"}, RunOpts{
		StdinReader: strings.NewReader(input),
	})

	expectedOut := `---
kind: Object
spec:
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
- image: index.docker.io/library/nginx@sha256:f7988fb6c02e0ce69257d9bd9cf37ae20a60f1df7563c3a2a6abe24160306b8d
`

	require.YAMLEq(t, expectedOut, out)
}

func TestResolveSuccessfulWithSearchRules(t *testing.T) {
	env := BuildEnv(t)
	kilnbase := Kilnbase{t, env.KilnbaseBinaryPath, Logger{}}

	input := `
kind: Object
spec:
- image: nginx:1.14.2
- sidecarImage: index.docker.io/library/nginx:1.14.2
- some_key: nginx:1.14.2
- nestedkey:
    data: |
      nested:
        image: nginx:1.14.2
- images:
