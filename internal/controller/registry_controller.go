/*
Copyright registry-operator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"slices"

	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	registryoperatordevv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
)

// RegistryReconciler reconciles a Registry object.
type RegistryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const registryFinalizer = "registry-operator.dev/finalizer"

// createInMemoryPod creates a pod with in-memory storage
// this functoin should be invoked through createPod function
// TODO: pod/deployment factory or builder.
func newInMemoryPodFromRegistry(registry *registryoperatordevv1alpha1.Registry) *apiv1.Pod {
	return &apiv1.Pod{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      registry.Name,
			Namespace: registry.Namespace,
			Labels: map[string]string{
				"app":      "registry",
				"registry": registry.Name,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  registry.Name,
					Image: "registry:2",
					Env: []apiv1.EnvVar{
						{
							Name:  "REGISTRY_STORAGE_INMEMORY",
							Value: "",
						},
						{
							Name:  "REGISTRY_STORAGE",
							Value: "",
						},
					},
				},
			},
		},
	}
}

// Supported storage types will be cases in this switch statement.
func newPodFromRegistry(registry *registryoperatordevv1alpha1.Registry) (*apiv1.Pod, error) {
	switch registry.Spec.StorageType {
	case registryoperatordevv1alpha1.StorageTypeInMemory:
		return newInMemoryPodFromRegistry(registry), nil
	default:
		return nil, fmt.Errorf("storage type %s not supported", registry.Spec.StorageType)
	}
}

// reconcileDeletion reconciles deletion of Registry CR
// For now it's just a pod deletion, but we will add more logic in the future.
func reconcileDeletion(ctx context.Context, r *RegistryReconciler, registry *registryoperatordevv1alpha1.Registry) error {
	l := log.FromContext(ctx)

	// Remove Pod
	pods := &apiv1.PodList{}
	err := r.Client.List(
		ctx,
		pods,
		client.InNamespace(registry.Namespace),
		client.MatchingLabels{"app": "registry", "registry": registry.Name},
	)
	if err != nil {
		l.Error(err, "Couldn't list pods", "namespace", registry.Namespace, "label", registry.Name)
		return err
	}

	idx := slices.IndexFunc(pods.Items, func(pod apiv1.Pod) bool {
		return pod.Name == registry.Name
	})

	if idx != -1 {
		l.Info("Deleting the pod")
		err := r.Client.Delete(ctx, &pods.Items[idx])
		if err != nil {
			l.Error(err, "Couldn't delete a pod", "name", registry.Name, "namespace", registry.Namespace)
			return err
		}
	}

	return nil
}

// reconcileFinalizers reconciles finalizers for Registry CR
// If we encounter more complex deletion logic in the future, this will need to be refactored.
func reconcileFinalizers(ctx context.Context, r *RegistryReconciler, registry *registryoperatordevv1alpha1.Registry) error {
	l := log.FromContext(ctx)

	if controllerutil.ContainsFinalizer(registry, registryFinalizer) { // Registry has finalizer
		if !registry.ObjectMeta.DeletionTimestamp.IsZero() { // Registry is being deleted if DeletionTimestamp is not zero
			if err := reconcileDeletion(ctx, r, registry); err != nil { // for now it's just a pod deletion, but we will add more logic in the future
				return err // we need to requeue if deletion failed
			}
			finalizersUpdated := controllerutil.RemoveFinalizer(registry, registryFinalizer)
			if finalizersUpdated {
				l.Info("Removing finalizer from Registry CR", "name", registry.Name, "namespace", registry.Namespace)
				if err := r.Update(ctx, registry); err != nil {
					l.Error(err, "Couldn't remove finalizer from Registry CR", "name", registry.Name, "namespace", registry.Namespace)
					return err
				}
			}
		}
		return nil // we don't need to do anything if finalizer is already present and Registry is not being deleted
	}

	if !registry.ObjectMeta.DeletionTimestamp.IsZero() { // Registry is being deleted and finalizer is not present
		return nil // nothing to do
	}

	// Add finalizer to Registry CR
	l.Info("Adding finalizer to Registry CR", "name", registry.Name, "namespace", registry.Namespace)
	_ = controllerutil.AddFinalizer(registry, registryFinalizer)
	if err := r.Update(ctx, registry); err != nil {
		l.Error(err, "Couldn't add finalizer to Registry CR", "name", registry.Name, "namespace", registry.Namespace)
		return err
	}

	return nil
}

// reconcileCreation reconciles creation of a pod for Registry CR.
func reconcileCreation(ctx context.Context, r *RegistryReconciler, registry *registryoperatordevv1alpha1.Registry) error {
	l := log.FromContext(ctx)

	// Check if pd for Registry CR already exists
	pod := &apiv1.Pod{}
	err := r.Client.Get(
		ctx,
		client.ObjectKey{
			Name:      registry.Name,
			Namespace: registry.Namespace,
		},
		pod,
	)
	if err != nil {
		// If pod doesn't exist, create a new one
		if apierrors.IsNotFound(err) {
			// Create a pod
			l.Info("Creating a pod")
			pod, err := newPodFromRegistry(registry)
			if err != nil {
				l.Error(err, "Couldn't create a pod", "name", registry.Name, "namespace", registry.Namespace)
				return err
			}

			if err := r.Client.Create(ctx, pod); err != nil {
				l.Error(err, "Couldn't create a pod", "name", registry.Name, "namespace", registry.Namespace)
				return err
			}
		} else { // If we got an error other than NotFound we need to requeue
			l.Error(err, "Couldn't get a pod", "name", registry.Name, "namespace", registry.Namespace)
			return err
		}
	}
	// If pod exists, do nothing for now
	return nil
}

//+kubebuilder:rbac:groups=registry-operator.dev,resources=registries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=registry-operator.dev,resources=registries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=registry-operator.dev,resources=registries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *RegistryReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Reconciling Registry", "name", request.Name, "namespace", request.Namespace)

	l.Info("Phase 1: Get the Registry CR from cluster")
	registry := &registryoperatordevv1alpha1.Registry{}
	if err := r.Get(ctx, request.NamespacedName, registry); err != nil {
		if apierrors.IsNotFound(err) {
			l.Info("Registry CR not found, it might have been deleted", "name", request.Name, "namespace", request.Namespace)
			return ctrl.Result{}, nil // don't requeue, nothing to do
		}

		l.Error(err, "Couldn't get Registry CR", "name", request.Name, "namespace", request.Namespace)
		return ctrl.Result{}, err // requeue
	}

	// If we got here, Registry CR exists

	l.Info("Phase 2: Reconcile finalizers, this includes deletion logic")
	if err := reconcileFinalizers(ctx, r, registry); err != nil {
		return ctrl.Result{}, err // requeue
	}

	l.Info("Phase 3: Create a pod if it doesn't exist")
	if registry.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := reconcileCreation(ctx, r, registry); err != nil {
			return ctrl.Result{}, err // requeue
		}
	} else {
		l.Info("Registry CR is being deleted, skipping pod creation")
	}

	l.Info("Phase 4: Done")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RegistryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryoperatordevv1alpha1.Registry{}).
		Complete(r)
}
