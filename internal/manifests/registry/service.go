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

	"github.com/registry-operator/registry-operator/internal/manifests"
	"github.com/registry-operator/registry-operator/internal/manifests/manifestutils"
	"github.com/registry-operator/registry-operator/internal/naming"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// headless and monitoring labels are to differentiate the base/headless/monitoring services from the clusterIP service.
const (
	serviceTypeLabel = "registry.registry-operator.dev/registry-service-type"
)

type ServiceType int

const (
	BaseServiceType ServiceType = iota
	HeadlessServiceType
)

func (s ServiceType) String() string {
	return [...]string{"base", "headless"}[s]
}

func convertServicePorts(containerPorts []corev1.ContainerPort) []corev1.ServicePort {
	servicePorts := []corev1.ServicePort{}

	for _, cp := range containerPorts {
		sp := corev1.ServicePort{
			Name:     cp.Name,
			Protocol: cp.Protocol,
			Port:     cp.ContainerPort,
		}

		if cp.Name != "" {
			sp.TargetPort = intstr.FromString(cp.Name)
		} else {
			sp.TargetPort = intstr.FromInt32(cp.ContainerPort)
		}

		servicePorts = append(servicePorts, sp)
	}

	return servicePorts
}

func Service(ctx context.Context, params manifests.Params) (*corev1.Service, error) {
	name := naming.Service(params.Registry.Name)
	labels := manifestutils.Labels(
		params.Registry.ObjectMeta,
		name,
		params.Registry.Spec.Image,
		ComponentRegistry,
		nil,
	)
	labels[serviceTypeLabel] = BaseServiceType.String()

	annotations, err := manifestutils.Annotations(params.Registry, nil)
	if err != nil {
		return nil, err
	}

	trafficPolicy := corev1.ServiceInternalTrafficPolicyCluster

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.Service(params.Registry.Name),
			Namespace:   params.Registry.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			InternalTrafficPolicy: &trafficPolicy,
			Selector:              manifestutils.SelectorLabels(params.Registry.ObjectMeta, ComponentRegistry),
			ClusterIP:             "",
			Ports:                 convertServicePorts(generateContainerPorts()),
		},
	}, nil
}
