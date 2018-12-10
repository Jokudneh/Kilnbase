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
	var path Path
	for _, part := range parts {
		switch typedPart := part.(type) {
		case string:
			path = append(path, NewPathPartFromString(typedPart))
		case int:
			path = append(path, NewPathPartFromIndex(typedPart))
		default:
			panic(fmt.Sprintf("Unexpected part type %T", typedPart))
		}
	}
	return path
}

func (p Path) AsStrings() []string {
	var result []string
	for _, part := range p {
		if part.MapKey == nil {
			panic(fmt.Sprintf("Unexpected non-map-key path part '%+v'", part))
		}
		result = append(result, *part.MapKey)
	}
	return result
}

func (p Path) AsString() string {
	var result []string
	for _, part := range p {
		result = append(result, part.AsString())
	}
	return strings.Join(result, ",")
