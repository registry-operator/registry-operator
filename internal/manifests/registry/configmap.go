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

package registry

import (
	"github.com/registry-operator/registry-operator/internal/manifests"
	"github.com/registry-operator/registry-operator/internal/manifests/manifestutils"
	"github.com/registry-operator/registry-operator/internal/naming"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml "sigs.k8s.io/yaml/goyaml.v2"
)

func ConfigMap(params manifests.Params) (*corev1.ConfigMap, error) {
	config := manifestutils.GenerateConfig(params.Registry.Spec)

	cfgYaml, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	hash, err := manifestutils.GetConfigMapSHA(config)
	if err != nil {
		return nil, err
	}

	name := naming.ConfigMap(params.Registry.Name, hash)
	labels := manifestutils.Labels(
		params.Registry.ObjectMeta,
		name,
		params.Registry.Spec.Image,
		ComponentRegistry,
		nil,
	)
	annotations, err := manifestutils.Annotations(params.Registry, nil)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   params.Registry.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: map[string]string{
			naming.DistributionConfig(): string(cfgYaml),
		},
	}, nil
}
