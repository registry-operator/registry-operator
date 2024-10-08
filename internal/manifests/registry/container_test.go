// Copyright Registry Operator contributors
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
	"github.com/registry-operator/registry-operator/internal/version"

	corev1 "k8s.io/api/core/v1"
)

func TestContainerDefault(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		Spec: registryv1alpha1.RegistrySpec{},
	}

	// test
	c := Container(registry)

	// verify
	assert.Equal(t, "distribution", c.Name)
	assert.Equal(t, version.GetRegistryImage(), c.Image)
	assert.Equal(t, c.ImagePullPolicy, corev1.PullIfNotPresent)
	assert.Empty(t, c.Resources)
	assert.Len(t, c.Ports, 1)
}

func TestContainerWithImageOverridden(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		Spec: registryv1alpha1.RegistrySpec{
			Image: "overridden-image",
		},
	}

	// test
	c := Container(registry)

	// verify
	assert.Equal(t, "overridden-image", c.Image)
}
