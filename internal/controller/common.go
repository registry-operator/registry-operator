// Copyright Registry Operator contributors
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
	"errors"
	"fmt"

	"github.com/registry-operator/registry-operator/internal/manifests"
	"github.com/registry-operator/registry-operator/internal/manifests/registry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// BuildRegistry returns the generation and collected errors of all manifests for a given instance.
func BuildRegistry(params manifests.Params) ([]client.Object, error) {
	builders := []manifests.Builder[manifests.Params]{
		registry.Build,
	}

	var resources []client.Object
	for _, builder := range builders {
		objs, err := builder(params)
		if err != nil {
			return nil, err
		}
		resources = append(resources, objs...)
	}

	return resources, nil
}

// reconcileDesiredObjects runs the reconcile process using the mutateFn over the given list of objects.
func reconcileDesiredObjects(
	ctx context.Context,
	kubeClient client.Client,
	owner metav1.Object,
	scheme *runtime.Scheme,
	desiredObjects []client.Object,
	ownedObjects map[types.UID]client.Object,
) error {
	log := log.FromContext(ctx)

	var errs []error
	for _, desired := range desiredObjects {
		l := log.WithValues(
			"object_name", desired.GetName(),
			"object_kind", desired.GetObjectKind(),
		)

		if setErr := ctrl.SetControllerReference(owner, desired, scheme); setErr != nil {
			l.Error(setErr, "failed to set controller owner reference to desired")
			errs = append(errs, setErr)
			continue
		}

		// existing is an object the controller runtime will hydrate for us
		// we obtain the existing object by deep copying the desired object because it's the most convenient way
		existing := desired.DeepCopyObject().(client.Object)
		mutateFn := manifests.MutateFuncFor(existing, desired)
		var op controllerutil.OperationResult
		crudErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			result, createOrUpdateErr := ctrl.CreateOrUpdate(ctx, kubeClient, existing, mutateFn)
			op = result
			return createOrUpdateErr
		})

		if crudErr != nil && errors.Is(crudErr, manifests.ErrImmutableChange) {
			l.Error(
				crudErr,
				"Detected immutable field change, trying to delete, new object will be created on next reconcile",
				"existing", klog.KObj(existing),
			)

			delErr := kubeClient.Delete(ctx, existing)
			if delErr != nil {
				return delErr
			}
			continue
		} else if crudErr != nil {
			l.Error(crudErr, "Failed to configure desired")
			errs = append(errs, crudErr)
			continue
		}

		l.V(1).Info(fmt.Sprintf("Desired has been %s", op))
		// This object is still managed by the operator, remove it from the list of objects to prune
		delete(ownedObjects, existing.GetUID())
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to create objects for %s: %w", owner.GetName(), errors.Join(errs...))
	}

	// Pruning owned objects in the cluster which are not should not be present after the reconciliation.
	err := deleteObjects(ctx, kubeClient, ownedObjects)
	if err != nil {
		return fmt.Errorf("failed to prune objects for %s: %w", owner.GetName(), err)
	}

	return nil
}

// getList queries the Kubernetes API to list the requested resource, setting the list l of type T.
func getList[T client.Object](
	ctx context.Context,
	cl client.Client,
	l T,
	options ...client.ListOption,
) (map[types.UID]client.Object, error) {
	ownedObjects := map[types.UID]client.Object{}
	list := &unstructured.UnstructuredList{}
	gvk, err := apiutil.GVKForObject(l, cl.Scheme())
	if err != nil {
		return nil, err
	}
	list.SetGroupVersionKind(gvk)
	err = cl.List(ctx, list, options...)
	if err != nil {
		return ownedObjects, fmt.Errorf("error listing %T: %w", l, err)
	}
	for i := range list.Items {
		ownedObjects[list.Items[i].GetUID()] = &list.Items[i]
	}
	return ownedObjects, nil
}

func deleteObjects(ctx context.Context, kubeClient client.Client, objects map[types.UID]client.Object) error {
	log := log.FromContext(ctx)

	// Pruning owned objects in the cluster which are not should not be present after the reconciliation.
	pruneErrs := []error{}
	for _, obj := range objects {
		l := log.WithValues(
			"object_name", obj.GetName(),
			"object_kind", obj.GetObjectKind().GroupVersionKind(),
		)

		l.Info("Pruning unmanaged resource")
		err := kubeClient.Delete(ctx, obj)
		if err != nil {
			l.Error(err, "Failed to delete resource")
			pruneErrs = append(pruneErrs, err)
		}
	}
	return errors.Join(pruneErrs...)
}
