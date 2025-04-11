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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/manifests"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeploymentNewDefault(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-instance",
			Namespace: "my-namespace",
		},
	}

	params := manifests.Params{
		Registry: registry,
	}

	// test
	d, err := Deployment(t.Context(), params)
	require.NoError(t, err)

	// verify
	assert.Equal(t, "my-instance-registry", d.Name)
	assert.Equal(t, "my-instance-registry", d.Labels["app.kubernetes.io/name"])

	assert.Len(t, d.Spec.Template.Spec.Containers, 1)

	expectedLabels := map[string]string{
		"app.kubernetes.io/component":  "registry",
		"app.kubernetes.io/instance":   "my-namespace.my-instance",
		"app.kubernetes.io/managed-by": "registry-operator",
		"app.kubernetes.io/name":       "my-instance-registry",
		"app.kubernetes.io/part-of":    "registry",
		"app.kubernetes.io/version":    "latest",
	}
	assert.Equal(t, expectedLabels, d.Spec.Template.Labels)

	expectedSelectorLabels := map[string]string{
		"app.kubernetes.io/component":  "registry",
		"app.kubernetes.io/instance":   "my-namespace.my-instance",
		"app.kubernetes.io/managed-by": "registry-operator",
		"app.kubernetes.io/part-of":    "registry",
	}
	assert.Equal(t, expectedSelectorLabels, d.Spec.Selector.MatchLabels)

	// the pod selector must be contained within pod spec's labels
	for k, v := range d.Spec.Selector.MatchLabels {
		assert.Equal(t, v, d.Spec.Template.Labels[k])
	}
}
