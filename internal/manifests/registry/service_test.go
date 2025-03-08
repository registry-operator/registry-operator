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
	"github.com/registry-operator/registry-operator/internal/naming"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDesiredService(t *testing.T) {
	t.Run("should return default service", func(t *testing.T) {
		params := manifests.Params{
			Registry: registryv1alpha1.Registry{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-instance",
					Namespace: "my-namespace",
				},
			},
		}

		actual, err := Service(params)
		require.NoError(t, err)
		assert.Equal(t, "my-instance-registry", actual.Name)
		assert.Equal(t, naming.Service(params.Registry.Name), actual.Name)
		assert.Len(t, actual.Spec.Ports, 1)
	})
}
