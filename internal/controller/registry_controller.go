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
	"slices"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	registryoperatordevv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
)

// RegistryReconciler reconciles a Registry object.
type RegistryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=registry-operator.dev,resources=registries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=registry-operator.dev,resources=registries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=registry-operator.dev,resources=registries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Registry object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *RegistryReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Hello from reconcile!")

	registry := &registryoperatordevv1alpha1.Registry{}
	err := r.Client.Get(ctx, request.NamespacedName, registry)
	if err != nil {
		l.Info("couldn't get registry CR deleting pod if exists")
		pods := &apiv1.PodList{}
		err = r.Client.List(ctx, pods, client.InNamespace(request.Namespace), client.MatchingLabels{"app": "registry", "registry": registry.Name})
		if err != nil {
			l.Error(err, "couldn't list pods")
			return ctrl.Result{}, err
		}

		idx := slices.IndexFunc(pods.Items, func(pod apiv1.Pod) bool {
			return pod.Name == registry.Name
		})

		if idx != -1 {
			err := r.Client.Delete(ctx, &pods.Items[idx])
			if err != nil {
				l.Error(err, "couldn't delete a pod")
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	}

	pods := &apiv1.PodList{}
	err = r.Client.List(ctx, pods, client.InNamespace(request.Namespace), client.MatchingLabels{"app": "registry", "registry": registry.Name})
	if err != nil {
		l.Error(err, "couldn't list pods")
		return ctrl.Result{}, err
	}

	// Check if pod list contains pod with the same name as registry CR
	idx := slices.IndexFunc(pods.Items, func(pod apiv1.Pod) bool {
		return pod.Name == registry.Name
	})

	if idx == -1 {
		// Create a new pod, mount secret from registry CR to the pod and start registry container
		l.Info("Creating a new pod")
		pod := &apiv1.Pod{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      registry.Name,
				Namespace: request.Namespace,
				Labels: map[string]string{
					"app":      "registry",
					"registry": registry.Name,
				},
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  request.Name,
						Image: "registry:2",
						VolumeMounts: []apiv1.VolumeMount{
							{
								Name:      "registry-secret",
								MountPath: "/var/lib/registry",
							},
						},
					},
				},
				Volumes: []apiv1.Volume{
					{
						Name: "registry-secret",
						VolumeSource: apiv1.VolumeSource{
							Secret: &apiv1.SecretVolumeSource{
								SecretName: registry.Spec.BucketAccessSecretName,
							},
						},
					},
				},
			},
		}

		if err := r.Client.Create(ctx, pod); err != nil {
			l.Error(err, "couldn't create a pod")
			return ctrl.Result{}, err
		}

		l.Info("Pod created")
		return ctrl.Result{}, nil
	}

	// Delete the pod if registry CR is deleted
	if !registry.ObjectMeta.DeletionTimestamp.IsZero() {
		l.Info("Deleting the pod")
		err := r.Client.Delete(ctx, &pods.Items[idx])
		if err != nil {
			l.Error(err, "couldn't delete a pod")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RegistryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryoperatordevv1alpha1.Registry{}).
		Complete(r)
}
