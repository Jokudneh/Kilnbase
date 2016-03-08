// Copyright 2024 The Kilnbase Authors.
// SPDX-License-Identifier: Apache-2.0

package resources

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

type Resource interface {
	GroupVersionResource() schema.GroupVersionResource
	Kind() string
	APIVersion() string
	APIGroup() string

	Name() string
	Description() string

	Annotations() map[string]string
	Labels() map[string]string

	Equal(res Resource) bool
	DeepCopy() Resource
	DeepCopyRaw() map[string]interface{}
	AsYAMLBytes() ([]byte, error)

	unstructured() unstructured.Unstructured     // private
	unstructuredPtr() *unstructured.Unstructured // private
}

type ResourceImpl struct {
	un        unstructured.Unstructured
	gvr       schema.GroupVersionResource
	transient bool
}

var _ Resource = &ResourceImpl{}

func NewResourceUnstructured(un unstructured.Unstructured, gvr schema.GroupVersionResource) *ResourceImpl {
	return &ResourceImpl{un: un, gvr: gvr}
}

func NewResourceFromBytes(data []byte) (*ResourceImpl, error) {
	var content map[string]interface{}

	err := yaml.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, nil
