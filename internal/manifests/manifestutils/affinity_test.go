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

package manifestutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
)

func TestAffinity(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// prepare
		registry := registryv1alpha1.Registry{
			Spec: registryv1alpha1.RegistrySpec{
				Affinity: nil,
			},
		}

		// test
		affinty := Affinity(registry)

		// verify
		assert.NotNil(t, affinty)
	})

	t.Run("affinity", func(t *testing.T) {
		// prepare
		expected := buildAffinity([]corev1.NodeSelectorRequirement{
			{
				Key:      "test.bar.io",
				Operator: corev1.NodeSelectorOpIn,
				Values:   []string{"foo"},
			},
		})

		registry := registryv1alpha1.Registry{
			Spec: registryv1alpha1.RegistrySpec{
				Affinity: expected,
			},
		}

		// test
		affinty := Affinity(registry)

		// verify
		assert.NotNil(t, affinty)
		assert.Equal(t, expected, affinty)
	})
}

func buildAffinity(nodeSelectors []corev1.NodeSelectorRequirement) *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: nodeSelectors,
					},
				},
			},
		},
	}
}
