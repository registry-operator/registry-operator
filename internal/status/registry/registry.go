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
	"context"
	"fmt"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/naming"
	"github.com/registry-operator/registry-operator/internal/version"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateRegistryStatus(ctx context.Context, cli client.Client, changed *registryv1alpha1.Registry) error {
	if changed.Status.Version == "" {
		// a version is not set, otherwise let the upgrade mechanism take care of it!
		changed.Status.Version = version.Registry()
	}

	// Set the scale replicas
	objKey := client.ObjectKey{
		Namespace: changed.GetNamespace(),
		Name:      naming.Registry(changed.Name),
	}

	var replicas int32
	var readyReplicas int32
	var statusImage string

	obj := &appsv1.Deployment{}
	if err := cli.Get(ctx, objKey, obj); err != nil {
		return fmt.Errorf("failed to get deployment status.replicas: %w", err)
	}
	replicas = obj.Status.Replicas
	readyReplicas = obj.Status.ReadyReplicas
	statusImage = obj.Spec.Template.Spec.Containers[0].Image

	changed.Status.Image = statusImage
	if replicas != 0 && replicas == readyReplicas {
		changed.Status.Ready = true
	}

	return nil
}
