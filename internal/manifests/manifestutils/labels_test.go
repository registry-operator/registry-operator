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
// Copyright The registry Authors

package manifestutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	collectorName      = "my-instance"
	collectorNamespace = "my-ns"
	taname             = "my-instance"
	tanamespace        = "my-ns"
)

func TestLabelsCommonSet(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      collectorName,
			Namespace: collectorNamespace,
		},
		Spec: registryv1alpha1.RegistrySpec{
			Image: "docker.io/library/registry:0.47.0",
		},
	}

	// test
	labels := Labels(registry.ObjectMeta, collectorName, registry.Spec.Image, "registry", []string{})
	assert.Equal(t, "registry-operator", labels["app.kubernetes.io/managed-by"])
	assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
	assert.Equal(t, "0.47.0", labels["app.kubernetes.io/version"])
	assert.Equal(t, "registry", labels["app.kubernetes.io/part-of"])
	assert.Equal(t, "registry", labels["app.kubernetes.io/component"])
}

func TestLabelsSha256Set(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      collectorName,
			Namespace: collectorNamespace,
		},
		Spec: registryv1alpha1.RegistrySpec{
			Image: "docker.io/library/registry@sha256:ac0192b549007e22998eb74e8d8488dcfe70f1489520c3b144a6047ac5efbe90",
		},
	}

	// test
	labels := Labels(registry.ObjectMeta, collectorName, registry.Spec.Image, "registry", []string{})
	assert.Equal(t, "registry-operator", labels["app.kubernetes.io/managed-by"])
	assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
	assert.Equal(t, "ac0192b549007e22998eb74e8d8488dcfe70f1489520c3b144a6047ac5efbe9", labels["app.kubernetes.io/version"])
	assert.Equal(t, "registry", labels["app.kubernetes.io/part-of"])
	assert.Equal(t, "registry", labels["app.kubernetes.io/component"])

	// prepare
	registryTag := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      collectorName,
			Namespace: collectorNamespace,
		},
		Spec: registryv1alpha1.RegistrySpec{
			Image: "docker.io/library/registry:2.8.3@sha256:ac0192b549007e22998eb74e8d8488dcfe70f1489520c3b144a6047ac5efbe90",
		},
	}

	// test
	labelsTag := Labels(registryTag.ObjectMeta, collectorName, registryTag.Spec.Image, "registry", []string{})
	assert.Equal(t, "registry-operator", labelsTag["app.kubernetes.io/managed-by"])
	assert.Equal(t, "my-ns.my-instance", labelsTag["app.kubernetes.io/instance"])
	assert.Equal(t, "2.8.3", labelsTag["app.kubernetes.io/version"])
	assert.Equal(t, "registry", labelsTag["app.kubernetes.io/part-of"])
	assert.Equal(t, "registry", labelsTag["app.kubernetes.io/component"])
}

func TestLabelsTagUnset(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      collectorName,
			Namespace: collectorNamespace,
		},
		Spec: registryv1alpha1.RegistrySpec{
			Image: "docker.io/library/registry",
		},
	}

	// test
	labels := Labels(registry.ObjectMeta, collectorName, registry.Spec.Image, "registry", []string{})
	assert.Equal(t, "registry-operator", labels["app.kubernetes.io/managed-by"])
	assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
	assert.Equal(t, "latest", labels["app.kubernetes.io/version"])
	assert.Equal(t, "registry", labels["app.kubernetes.io/part-of"])
	assert.Equal(t, "registry", labels["app.kubernetes.io/component"])
}

func TestLabelsPropagateDown(t *testing.T) {
	// prepare
	registry := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"myapp":                  "mycomponent",
				"app.kubernetes.io/name": "test",
			},
		},
		Spec: registryv1alpha1.RegistrySpec{
			Image: "docker.io/library/registry",
		},
	}

	// test
	labels := Labels(registry.ObjectMeta, collectorName, registry.Spec.Image, "registry", []string{})

	// verify
	assert.Len(t, labels, 7)
	assert.Equal(t, "mycomponent", labels["myapp"])
	assert.Equal(t, "test", labels["app.kubernetes.io/name"])
}

func TestLabelsFilter(t *testing.T) {
	registry := registryv1alpha1.Registry{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"test.bar.io": "foo", "test.foo.io": "bar"},
		},
	}

	// This requires the filter to be in regex match form and not the other simpler wildcard one.
	labels := Labels(registry.ObjectMeta, collectorName, "latest", "registry", []string{".*.bar.io"})

	// verify
	assert.Len(t, labels, 7)
	assert.NotContains(t, labels, "test.bar.io")
	assert.Equal(t, "bar", labels["test.foo.io"])
}
