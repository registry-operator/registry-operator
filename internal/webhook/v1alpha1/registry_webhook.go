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

package v1alpha1

import (
	"context"
	"fmt"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/webhook/validation"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupRegistryWebhookWithManager registers the webhook for Registry in the manager.
func SetupRegistryWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&registryv1alpha1.Registry{}).
		WithValidator(&RegistryCustomValidator{}).
		WithDefaulter(&RegistryCustomDefaulter{}).
		Complete()
}

//nolint:lll // long URLs
// +kubebuilder:webhook:path=/mutate-registry-operator-dev-v1alpha1-registry,mutating=true,failurePolicy=fail,sideEffects=None,groups=registry-operator.dev,resources=registries,verbs=create;update,versions=v1alpha1,name=mregistry-v1alpha1.kb.io,admissionReviewVersions=v1

// RegistryCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Registry when those are created or updated.
type RegistryCustomDefaulter struct{}

var _ webhook.CustomDefaulter = &RegistryCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Registry.
func (d *RegistryCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	log := ctrl.LoggerFrom(ctx).WithName("registry-resource")

	registry, ok := obj.(*registryv1alpha1.Registry)
	if !ok {
		return fmt.Errorf("expected an Registry object but got %T", obj)
	}

	log.V(4).Info("Defaulting for Registry", "name", registry.GetName())

	if validation.PopulatedFields(registry.Spec.Storage) == 0 {
		registry.Spec.Storage = registryv1alpha1.Storage{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				SizeLimit: ptr.To(resource.MustParse("200Mi")),
			},
		}
	}

	htpasswd := registry.Spec.Auth.Htpasswd
	if htpasswd != nil &&
		htpasswd.Secret.Optional != nil &&
		htpasswd.Secret.Key == "" {
		htpasswd.Secret.Key = "auth"
	}

	return nil
}

//nolint:lll // kubebuilder directives
// +kubebuilder:webhook:path=/validate-registry-operator-dev-v1alpha1-registry,mutating=false,failurePolicy=fail,sideEffects=None,groups=registry-operator.dev,resources=registries,verbs=create;update,versions=v1alpha1,name=vregistry-v1alpha1.kb.io,admissionReviewVersions=v1

// RegistryCustomValidator struct is responsible for validating the Registry resource
// when it is created, updated, or deleted.
type RegistryCustomValidator struct{}

var _ webhook.CustomValidator = &RegistryCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Registry.
func (v *RegistryCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	log := ctrl.LoggerFrom(ctx).WithName("registry-resource")

	registry, ok := obj.(*registryv1alpha1.Registry)
	if !ok {
		return nil, fmt.Errorf("expected a Registry object but got %T", obj)
	}

	log.V(3).Info("Validation for Registry upon creation", "name", registry.GetName())

	return v.warn(registry), v.validate(registry)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Registry.
func (v *RegistryCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	log := ctrl.LoggerFrom(ctx).WithName("registry-resource")

	newRegistry, ok := newObj.(*registryv1alpha1.Registry)
	if !ok {
		return nil, fmt.Errorf("expected a Registry object for the newObj but got %T", newObj)
	}

	log.V(3).Info("Validation for Registry upon update", "name", newRegistry.GetName())

	return v.warn(newRegistry), v.validate(newRegistry)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Registry.
func (v *RegistryCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	log := ctrl.LoggerFrom(ctx).WithName("registry-resource")

	registry, ok := obj.(*registryv1alpha1.Registry)
	if !ok {
		return nil, fmt.Errorf("expected a Registry object for the newObj but got %T", obj)
	}

	log.V(4).Info("Validation for Registry upon update", "name", registry.GetName())

	return nil, nil
}

func (v *RegistryCustomValidator) validate(registry *registryv1alpha1.Registry) error {
	var allErrs field.ErrorList

	if !validation.HasAtMostOne(registry.Spec.Storage) {
		err := field.Invalid(
			field.NewPath("spec").Child("storage"),
			registry.Spec.Storage,
			"must contain at most one value",
		)
		allErrs = append(allErrs, err)
	}

	if len(allErrs) != 0 {
		return apierrors.NewInvalid(
			schema.GroupKind{
				Group: registryv1alpha1.GroupVersion.Group,
				Kind:  registryv1alpha1.RegistryKind,
			},
			registry.Name,
			allErrs,
		)
	}
	return nil
}

func (v *RegistryCustomValidator) warn(registry *registryv1alpha1.Registry) admission.Warnings {
	var warns admission.Warnings

	if registry.Spec.Replicas > 1 {
		warns = append(warns,
			"If replicas > 1 and file/block storage is used, there is no data consistency between Registry replicas.",
		)
	}

	htpasswd := registry.Spec.Auth.Htpasswd
	if htpasswd != nil && htpasswd.Secret.Optional != nil {
		warns = append(warns,
			"Htpasswd optional setting is ignored.",
		)
	}

	return warns
}
