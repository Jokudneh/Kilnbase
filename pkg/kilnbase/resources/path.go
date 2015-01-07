// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Path []*PathPart

type PathPart struct {
	MapKey     *string             `json:"mapKey,omitempty"`
	ArrayIndex *PathPartArrayIndex `json:"arrayIndex,omitempty"`
}

var _ json.Unmarshaler = &PathPart{}
var _ json.Marshaler = Path{}

type PathPartArrayIndex struct {
	Index *int  `json:"index,omitempty"`
	All   *bool `json:"allIndexes,omitempty"`
}

func NewPathFromStrings(strs []string) Path {
	var path Path
	for _, str := range strs {
		path = append(path, NewPathPartFromString(str))
	}
	return path
}

func NewPathFromInterfaces(parts []interface{}) Path {
