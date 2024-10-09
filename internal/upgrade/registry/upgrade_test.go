// Copyright The Registry Operator Authors
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

package registry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/upgrade/registry"
	"github.com/registry-operator/registry-operator/internal/version"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
)

func TestVersionsShouldNotBeChanged(t *testing.T) {
	nsn := types.NamespacedName{Name: "my-instance", Namespace: "default"}
	for _, tt := range []struct {
		desc            string
		v               string
		expectedV       string
		failureExpected bool
	}{
		{"new-instance", "", "", false},
		{"newer-than-our-newest", "100.0.0", "100.0.0", false},
		{"unparseable", "unparseable", "unparseable", true},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			// prepare
			existing := makeRegistry(nsn)
			existing.Status.Version = tt.v

			currentV := version.Get()
			currentV.Registry = registry.Latest.String()

			up := &registry.VersionUpgrade{
				Version:  currentV,
				Client:   k8sClient,
				Recorder: record.NewFakeRecorder(registry.RecordBufferSize),
			}

			// test
			res, err := up.ManagedInstance(context.Background(), existing)
			if tt.failureExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// verify
			assert.Equal(t, tt.expectedV, res.Status.Version)
		})
	}
}

func makeRegistry(nsn types.NamespacedName) registryv1alpha1.Registry {
	return registryv1alpha1.Registry{
		Spec: registryv1alpha1.RegistrySpec{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      nsn.Name,
			Namespace: nsn.Namespace,
		},
	}
}
