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
	"path"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/naming"
	"github.com/registry-operator/registry-operator/internal/version"

	corev1 "k8s.io/api/core/v1"
)

const (
	distributionPortDefault = 5000
	metricsPortDefault      = 5001
	configMountPath         = "/etc/distribution"
	storageMountPath        = "/var/lib/registry"
)

func generateContainerPorts() []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          naming.RegistryDistributionPort(),
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: distributionPortDefault,
		},
		{
			Name:          naming.RegistryMetricsPort(),
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: metricsPortDefault,
		},
	}
}

func generateVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      naming.ConfigVolume(),
			ReadOnly:  true,
			MountPath: configMountPath,
		},
		{
			Name:      naming.StorageVolume(),
			ReadOnly:  false,
			MountPath: storageMountPath,
		},
	}
}

func generateResources(res *corev1.ResourceRequirements) corev1.ResourceRequirements {
	if res != nil {
		return *res
	}
	return corev1.ResourceRequirements{}
}

// Container builds a container for the given registry.
func Container(registry registryv1alpha1.Registry) corev1.Container {
	image := registry.Spec.Image
	if len(image) == 0 {
		image = version.GetRegistryImage()
	}

	return corev1.Container{
		Name:            naming.Container(),
		Image:           image,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Command:         []string{"registry"},
		Args:            []string{"serve", path.Join(configMountPath, naming.DistributionConfig())},
		Ports:           generateContainerPorts(),
		VolumeMounts:    generateVolumeMounts(),
		Resources:       generateResources(registry.Spec.Resources),
	}
}
