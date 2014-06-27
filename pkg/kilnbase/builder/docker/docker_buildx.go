// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	ctlb "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	regname "github.com/google/go-containerregistry/pkg/name"
)

/*

Useful docs:
- https://cloudolife.com/2022/03/05/Infrastructure-as-Code-IaC/Container/Docker/Docker-buildx-support-multiple-architectures-images/

Commands:
- docker buildx create --name mybuilder
- docker buildx use mybuilder
- docker buildx inspect --bootstrap
or
- docker buildx create --name mybuilder --use

Example:

$ docker buildx build . -t k14stest/test-image:buildx-test-img --load --progress=plain
...snip...
#10 exporting to oci image format
#10 exporting layers done
#10 exporting manifest sha256:cd6d662fcf06854810c83e4bff197ef52e65b39c13baa143107ebd32cc6f8636 0.0s done
#10 exporting config sha256:e3bdd21522d99c37f355c70bef342c99cb1e19ee4cf8014001d750a73a5b8d42 done
#10 sending tarball
#10 sending tarball 0.1s done
#10 DONE 0.2s
#11 importing to docker
#11 DONE 0.0s

$ docker buildx build . -t k14stest/test-image:buildx-test-img --push --progress=plain
...snip...
#10 [auth] k14stest/test-image:pull,push token for registry-1.docker.io
#10 DONE 0.0s
#11 exporting to image
#11 exporting layers done
#11 exporting manifest sha256:cd6d662fcf06854810c83e4bff197ef52e65b39c13baa143107ebd32cc6f8636 done
#11 exporting config sha256:e3bdd21522d99c37f355c70bef342c99cb1e19ee4cf8014001d750a73a5b8d42 done
#11 pushing layers
#11 pushing layers 0.3s done
#11 pushing manifest for docker.io/k14stest/test-image:buildx-test-img@sha256:cd6d662fcf06854810c83e4bff197ef52e65b39c13baa143107ebd32cc6f8636
#11 pushing manifest for docker.io/k14stest/test-image:buildx-test-img@sha256:cd6d662fcf06854810c83e4bff197ef52e65b39c13baa143107ebd32cc6f8636 0.2s done
#11 DONE 0.6s

$ docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 . --push -t k14stest/test-image:buildx-test-img --progress=plain
...snip...
