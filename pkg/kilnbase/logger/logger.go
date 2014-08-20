// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

const noBytesWritten = 0

type Logger struct {
	writer     io.Writer
	writerLock *sync.Mutex
}
