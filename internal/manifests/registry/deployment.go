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
	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/manifests"
	"github.com/registry-operator/registry-operator/internal/manifests/manifestutils"
	"github.com/registry-operator/registry-operator/internal/naming"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func generateConfigVolume(registry, hash string) corev1.Volume {
	config := corev1.KeyToPath{
		Key:  naming.DistributionConfig(),
		Path: naming.DistributionConfig(),
	}

	return corev1.Volume{
		Name: naming.ConfigVolume(),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				Items: []corev1.KeyToPath{
					config,
				},
				LocalObjectReference: corev1.LocalObjectReference{
					Name: naming.ConfigMap(registry, hash),
				},
			},
		},
	}
}

func generateStorageVolume(registry registryv1alpha1.Registry) corev1.Volume {
	source := corev1.VolumeSource{
		EmptyDir: &corev1.EmptyDirVolumeSource{},
	}

	storage := registry.Spec.Storage
	if storage.EmptyDir != nil {
		source = corev1.VolumeSource{
			EmptyDir: storage.EmptyDir,
		}
	} else if storage.Ephemeral != nil {
		source = corev1.VolumeSource{
			Ephemeral: storage.Ephemeral,
		}
	} else if storage.HostPath != nil {
		source = corev1.VolumeSource{
			HostPath: storage.HostPath,
		}
	} else if storage.PersistentVolumeClaim != nil {
		source = corev1.VolumeSource{
			PersistentVolumeClaim: storage.PersistentVolumeClaim,
		}
	} else if storage.PersistentVolumeClaimTemplate != nil {
		source = corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: naming.PersistentVolumeClaim(registry.Name),
				ReadOnly:  false,
			},
		}
	}

	return corev1.Volume{
		Name:         naming.StorageVolume(),
		VolumeSource: source,
	}
}

// Deployment builds the deployment for the given instance.
func Deployment(params manifests.Params) (*appsv1.Deployment, error) {
	name := naming.Registry(params.Registry.Name)
	labels := manifestutils.Labels(
		params.Registry.ObjectMeta,
		name,
		params.Registry.Spec.Image,
		ComponentRegistry,
		nil,
	)
	annotations, err := manifestutils.Annotations(params.Registry, nil)
	if err != nil {
		return nil, err
	}

	podAnnotations, err := manifestutils.PodAnnotations(params.Registry, nil)
	if err != nil {
		return nil, err
	}

	hash, err := manifestutils.GetConfigMapSHA(manifestutils.GenerateConfig(params.Registry.Spec))
	if err != nil {
		return nil, err
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   params.Registry.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(params.Registry.Spec.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: manifestutils.SelectorLabels(params.Registry.ObjectMeta, ComponentRegistry),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: podAnnotations,
				},
				Spec: corev1.PodSpec{
					Affinity: manifestutils.Affinity(params.Registry),
					Containers: []corev1.Container{
						Container(params.Registry),
					},
					Volumes: []corev1.Volume{
						generateConfigVolume(params.Registry.Name, hash),
						generateStorageVolume(params.Registry),
					},
				},
			},
		},
	}, nil
}
