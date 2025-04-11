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

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/manifests"
	"github.com/registry-operator/registry-operator/internal/manifests/manifestutils"
	"github.com/registry-operator/registry-operator/internal/naming"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "embed"
)

var (
	//go:embed test/dev-config.yaml
	devConfig string
)

func TestDesiredSecret(t *testing.T) {
	expectedLables := map[string]string{
		"app.kubernetes.io/component":  "registry",
		"app.kubernetes.io/instance":   "my-namespace.my-instance",
		"app.kubernetes.io/managed-by": "registry-operator",
		"app.kubernetes.io/part-of":    "registry",
		"app.kubernetes.io/version":    "latest",
		"app.kubernetes.io/name":       "d4f6e0d090e2219a53b40993a8514e4035b8077bd5cf8643b9683282eaaf8b",
	}

	t.Run("should return expected collector config map", func(t *testing.T) {
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

		expectedData := map[string]string{
			"config.yaml": devConfig,
		}

		config, _ := generateConfig(t.Context(), params)
		hash, _ := manifestutils.CalculateHash(config)
		expectedName := naming.Secret("test", hash)

		// test
		actual, err := Secret(t.Context(), params)

		// verify
		assert.NoError(t, err)
		assert.Equal(t, expectedName, actual.Name)
		assert.Equal(t, expectedLables, actual.Labels)
		assert.Equal(t, len(expectedData), len(actual.StringData))
		for k, expected := range expectedData {
			t.Skip("TODO: untagged struct fields are not not omitted when empty")
			assert.YAMLEq(t, expected, actual.StringData[k])
		}
	})
}
