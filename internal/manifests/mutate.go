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

package manifests

import (
	"fmt"
	"reflect"

	"dario.cat/mergo"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ImmutableFieldChangeErr struct {
	Field string
}

func (e *ImmutableFieldChangeErr) Error() string {
	return fmt.Sprintf("Immutable field change attempted: %s", e.Field)
}

var (
	ImmutableChangeErr *ImmutableFieldChangeErr
)

// MutateFuncFor returns a mutate function based on the
// existing resource's concrete type. It supports currently
// only the following types or else panics:
// - ConfigMap
// - Deployment
// - Service
// In order for the operator to reconcile other types, they must be added here.
// The function returned takes no arguments but instead uses the existing and desired inputs here. Existing is expected
// to be set by the controller-runtime package through a client get call.
func MutateFuncFor(existing, desired client.Object) controllerutil.MutateFn {
	return func() error {
		// Get the existing annotations and override any conflicts with the desired annotations
		// This will preserve any annotations on the existing set.
		existingAnnotations := existing.GetAnnotations()
		if err := mergeWithOverride(&existingAnnotations, desired.GetAnnotations()); err != nil {
			return err
		}
		existing.SetAnnotations(existingAnnotations)

		// Get the existing labels and override any conflicts with the desired labels
		// This will preserve any labels on the existing set.
		existingLabels := existing.GetLabels()
		if err := mergeWithOverride(&existingLabels, desired.GetLabels()); err != nil {
			return err
		}
		existing.SetLabels(existingLabels)

		if ownerRefs := desired.GetOwnerReferences(); len(ownerRefs) > 0 {
			existing.SetOwnerReferences(ownerRefs)
		}

		switch existing.(type) {
		case *appsv1.Deployment:
			dpl := existing.(*appsv1.Deployment)
			wantDpl := desired.(*appsv1.Deployment)
			return mutateDeployment(dpl, wantDpl)

		case *corev1.ConfigMap:
			cm := existing.(*corev1.ConfigMap)
			wantCm := desired.(*corev1.ConfigMap)
			mutateConfigMap(cm, wantCm)

		case *corev1.PersistentVolumeClaim:
			pvc := existing.(*corev1.PersistentVolumeClaim)
			wantPvc := desired.(*corev1.PersistentVolumeClaim)
			mutatePersistentVolumeClaim(pvc, wantPvc)

		case *corev1.Service:
			svc := existing.(*corev1.Service)
			wantSvc := desired.(*corev1.Service)
			mutateService(svc, wantSvc)

		default:
			t := reflect.TypeOf(existing).String()
			return fmt.Errorf("missing mutate implementation for resource type: %s", t)
		}
		return nil
	}
}

func mergeWithOverride(dst, src interface{}) error {
	return mergo.Merge(dst, src, mergo.WithOverride)
}

func mutatePersistentVolumeClaim(existing, desired *corev1.PersistentVolumeClaim) {
	existing.Spec.Resources.Requests = desired.Spec.Resources.Requests
}

func mutateService(existing, desired *corev1.Service) {
	existing.Spec.Ports = desired.Spec.Ports
	existing.Spec.Selector = desired.Spec.Selector
}

func mutateConfigMap(existing, desired *corev1.ConfigMap) {
	existing.BinaryData = desired.BinaryData
	existing.Data = desired.Data
}

func hasImmutableLabelChange(existingSelectorLabels, desiredLabels map[string]string) error {
	for k, v := range existingSelectorLabels {
		if vv, ok := desiredLabels[k]; !ok || vv != v {
			return &ImmutableFieldChangeErr{Field: "Spec.Template.Metadata.Labels"}
		}
	}
	return nil
}

func mutatePodTemplate(existing, desired *corev1.PodTemplateSpec) error {
	if err := mergeWithOverride(&existing.Labels, desired.Labels); err != nil {
		return err
	}

	if err := mergeWithOverride(&existing.Annotations, desired.Annotations); err != nil {
		return err
	}

	existing.Spec = desired.Spec

	return nil

}

func mutateDeployment(existing, desired *appsv1.Deployment) error {
	if !existing.CreationTimestamp.IsZero() {
		if !apiequality.Semantic.DeepEqual(desired.Spec.Selector, existing.Spec.Selector) {
			return &ImmutableFieldChangeErr{Field: "Spec.Selector"}
		}
		if err := hasImmutableLabelChange(existing.Spec.Selector.MatchLabels, desired.Spec.Template.Labels); err != nil {
			return err
		}
	}

	existing.Spec.MinReadySeconds = desired.Spec.MinReadySeconds
	existing.Spec.Paused = desired.Spec.Paused
	existing.Spec.ProgressDeadlineSeconds = desired.Spec.ProgressDeadlineSeconds
	existing.Spec.Replicas = desired.Spec.Replicas
	existing.Spec.RevisionHistoryLimit = desired.Spec.RevisionHistoryLimit
	existing.Spec.Strategy = desired.Spec.Strategy

	if err := mutatePodTemplate(&existing.Spec.Template, &desired.Spec.Template); err != nil {
		return err
	}

	return nil
}
