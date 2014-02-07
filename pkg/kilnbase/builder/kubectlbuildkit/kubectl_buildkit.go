// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"

	ctlb "github.com/Jokudneh/Kilnbase/pkg/kilnbase/builder"
	ctlconf "github.com/Jokudneh/Kilnbase/pkg/kilnbase/config"
	ctllog "github.com/Jokudneh/Kilnbase/pkg/kilnbase/logger"
	regname "github.com/google/go-containerregistry/pkg/name"
)

var (
	// Example output that includes final digest:
	//   ...
	// 	 #10 exporting layers done
	// 	 #10 exporting manifest sha256:55d863c4231ec285b88516942fd5d636216c36b6a686a20bf28d1aa5125c16b7 0.0s done
	// 	 #10 exporting config sha256:aa1fce99c57c864cc8d98f7ff4a54bc3b9b7e3f63a741de926a3393bc76f5986 0.0s done
	//   ...
	//   #10 exporting layers done
	//   #10 exporting manifest sha256:55d863c4231ec285b88516942fd5d636216c36b6a686a20bf28d1aa5125c16b7 done
	//   #10 exporting config sha256:aa1fce99c57c864cc8d98f7ff4a54bc3b9b7e3f63a741de926a3393bc76f5986 done
	//   #10 pushing layers
	kubectlBuildkitDigest = regexp.MustCompile("exporting manifest (sha256:)?([0-9a-z]+) ")
)

type KubectlBuildkit struct {
	logger ctllog.Logger
}

func NewKubectlBuildkit(logger ctllog.Logger) KubectlBuildkit {
	return KubectlBuildkit{logger}
}

func (d KubectlBuildkit) BuildAndPush(image, directory string,
	imgDst *ctlconf.ImageDestination, opts ctlconf.SourceKubectlBuildkitOpts) (string, error) {

	tagRef, err := d.tagRef(image, imgDst)
	if err != nil {
		return "", err
	}

	prefixedLogger := d.logger.NewPrefixedWriter(image + " | ")
