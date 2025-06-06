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

package registry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/upgrade/registry"
	"github.com/registry-operator/registry-operator/internal/version"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
)

func TestV3_0_0(t *testing.T) {
	registryInstance := registryv1alpha1.Registry{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Registry",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "registry-my-instance",
			Namespace: "somewhere",
		},
		Status: registryv1alpha1.RegistryStatus{
			Version: "3.0.0",
		},
		Spec: registryv1alpha1.RegistrySpec{},
	}

	versionUpgrade := &registry.VersionUpgrade{
		Version:  version.Get(),
		Client:   k8sClient,
		Recorder: record.NewFakeRecorder(registry.RecordBufferSize),
	}

	reg, err := versionUpgrade.ManagedInstance(context.Background(), registryInstance)
	assert.NoError(t, err)
	assert.Equal(t, registryInstance, reg)
}
