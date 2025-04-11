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

package controller

import (
	"context"
	"errors"
	"reflect"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var watchLogger = ctrl.Log.WithName("watch")

var errInvalidType = errors.New("invalid type")

func (r *RegistryReconciler) MapS3Secrets(ctx context.Context, obj client.Object) []reconcile.Request {
	if _, ok := obj.(*corev1.Secret); !ok {
		t := reflect.TypeOf(obj).String()
		watchLogger.Error(errInvalidType, "Invalid type of Object, expected Secret", "type", t)
		return nil
	}

	list := &registryv1alpha1.RegistryList{}
	if err := r.List(ctx, list, client.InNamespace(obj.GetNamespace())); err != nil {
		watchLogger.Error(err, "Failed to list Registries", "namespace", obj.GetNamespace())
		return nil
	}

	objects := map[types.UID]types.NamespacedName{}

	for _, reg := range list.Items {
		if s3 := reg.Spec.Storage.S3; s3 != nil {
			name := obj.GetName()

			refs := []*registryv1alpha1.SecretKeySelector{
				&s3.BucketName,
				&s3.Region,
				s3.AccessKey,
				s3.SecretKey,
				s3.EndpointURL,
			}

			for _, ref := range refs {
				if ref != nil && ref.Name == name {
					objects[reg.GetUID()] = types.NamespacedName{
						Name:      reg.GetName(),
						Namespace: reg.GetNamespace(),
					}
					break // only need one match
				}
			}
		}
	}

	reqs := []reconcile.Request{}
	for _, req := range objects {
		reqs = append(reqs, reconcile.Request{NamespacedName: req})
	}

	return reqs
}

var (
	s3SecretPredicate predicate.Funcs = predicate.Funcs{
		CreateFunc: func(_ event.TypedCreateEvent[client.Object]) bool { return true },
		DeleteFunc: func(_ event.TypedDeleteEvent[client.Object]) bool { return true },
		UpdateFunc: func(e event.TypedUpdateEvent[client.Object]) bool {
			new, ok := e.ObjectNew.(*corev1.Secret)
			if !ok {
				t := reflect.TypeOf(new).String()
				watchLogger.Error(errInvalidType, "Invalid type of New, expected Secret", "type", t)
				return false
			}
			old, ok := e.ObjectOld.(*corev1.Secret)
			if !ok {
				t := reflect.TypeOf(old).String()
				watchLogger.Error(errInvalidType, "Invalid type of Old, expected Secret", "type", t)
				return false
			}

			return !reflect.DeepEqual(
				new.Data,
				old.Data,
			)
		},
	}
)
