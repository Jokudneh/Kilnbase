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

func NewLogger(writer io.Writer) Logger {
	return Logger{writer: writer, writerLock: &sync.Mutex{}}
}

func (l Logger) NewPrefixedWriter(prefix string) *PrefixWriter {
	return &PrefixWriter{prefix, l.writer, l.writerLock}
}

type PrefixWriter struct {
	prefix     string
	writer     io.Writer
