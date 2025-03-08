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

package registry

import (
	"github.com/registry-operator/registry-operator/internal/manifests"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ComponentRegistry = "registry"
)

// Build creates the manifest for the registry resource.
func Build(params manifests.Params) ([]client.Object, error) {
	var resourceManifests []client.Object
	var manifestFactories []manifests.K8sManifestFactory[manifests.Params]

	manifestFactories = append(manifestFactories, []manifests.K8sManifestFactory[manifests.Params]{
		manifests.Factory(Deployment),
		manifests.Factory(ConfigMap),
		manifests.Factory(Service),
	}...)

	for _, factory := range manifestFactories {
		res, err := factory(params)
		if err != nil {
			return nil, err
		} else if manifests.ObjectIsNotNil(res) {
			resourceManifests = append(resourceManifests, res)
		}
	}

	return resourceManifests, nil
}
