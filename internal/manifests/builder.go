// Copyright 2025 The Registry Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Additional copyrights:
// Copyright The OpenTelemetry Authors

package manifests

import (
	"context"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Builder[Params any] func(ctx context.Context, params Params) ([]client.Object, error)

type ManifestFactory[T client.Object, Params any] func(ctx context.Context, params Params) (T, error)
type K8sManifestFactory[Params any] ManifestFactory[client.Object, Params]

func Factory[T client.Object, Params any](f ManifestFactory[T, Params]) K8sManifestFactory[Params] {
	return func(ctx context.Context, params Params) (client.Object, error) {
		return f(ctx, params)
	}
}

// ObjectIsNotNil ensures that we only create an object IFF it isn't nil,
// and it's concrete type isn't nil either. This works around the Go type system
// by using reflection to verify its concrete type isn't nil.
func ObjectIsNotNil(obj client.Object) bool {
	return obj != nil && !reflect.ValueOf(obj).IsNil()
}
