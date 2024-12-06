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

package controller

import (
	"context"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/manifests"
	"github.com/registry-operator/registry-operator/internal/manifests/manifestutils"
	"github.com/registry-operator/registry-operator/internal/manifests/registry"
	registrystatus "github.com/registry-operator/registry-operator/internal/status/registry"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	registryFinalizer = "registry.registry-operator.dev/finalizer"
)

// RegistryReconciler reconciles a Registry object.
type RegistryReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=registry-operator.dev,resources=registries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=registry-operator.dev,resources=registries/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=registry-operator.dev,resources=registries/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// SetupWithManager sets up the controller with the Manager.
func (r *RegistryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryv1alpha1.Registry{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *RegistryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx, "registry", klog.KRef(req.Namespace, req.Name))

	var instance registryv1alpha1.Registry
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "unable to fetch Registry")
		}

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	params, err := r.GetParams(instance)
	if err != nil {
		log.Error(err, "Failed to create manifest.Params")
		return ctrl.Result{}, err
	}

	// We have a deletion, short circuit and let the deletion happen
	if deletionTimestamp := instance.GetDeletionTimestamp(); deletionTimestamp != nil {
		if controllerutil.ContainsFinalizer(&instance, registryFinalizer) {
			// If the finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err = r.finalizeRegistry(ctx, params); err != nil {
				return ctrl.Result{}, err
			}

			// Once all finalizers have been
			// removed, the object will be deleted.
			if controllerutil.RemoveFinalizer(&instance, registryFinalizer) {
				err = r.Update(ctx, &instance)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}

		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(&instance, registryFinalizer) {
		if controllerutil.AddFinalizer(&instance, registryFinalizer) {
			err = r.Update(ctx, &instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	desiredObjects, buildErr := BuildRegistry(params)
	if buildErr != nil {
		return ctrl.Result{}, buildErr
	}

	ownedObjects, err := r.findRegistryOwnedObjects(ctx, params)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = reconcileDesiredObjects(ctx, r.Client, &instance, params.Scheme, desiredObjects, ownedObjects)
	return registrystatus.HandleReconcileStatus(ctx, params, instance, err)
}

func (r *RegistryReconciler) GetParams(instance registryv1alpha1.Registry) (manifests.Params, error) {
	p := manifests.Params{
		Client:   r.Client,
		Registry: instance,
		Scheme:   r.Scheme,
		Recorder: r.Recorder,
	}

	return p, nil
}

func (r *RegistryReconciler) findRegistryOwnedObjects(
	ctx context.Context,
	params manifests.Params,
) (map[types.UID]client.Object, error) {
	ownedObjects := map[types.UID]client.Object{}
	ownedObjectTypes := []client.Object{}

	listOps := &client.ListOptions{
		Namespace: params.Registry.Namespace,
		LabelSelector: labels.SelectorFromSet(
			manifestutils.SelectorLabels(params.Registry.ObjectMeta, registry.ComponentRegistry),
		),
	}
	for _, objectType := range ownedObjectTypes {
		objs, err := getList(ctx, r.Client, objectType, listOps)
		if err != nil {
			return nil, err
		}
		for uid, object := range objs {
			ownedObjects[uid] = object
		}
	}

	return ownedObjects, nil
}

func (r *RegistryReconciler) finalizeRegistry(_ context.Context, _ manifests.Params) error {
	return nil
}
