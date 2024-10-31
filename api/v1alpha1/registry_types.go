/*
Copyright The Registry Operator Authors

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RegistrySpec defines the desired state of Registry.
type RegistrySpec struct {
	// Image indicates the container image to use for the Registry.
	// +optional
	Image string `json:"image,omitempty"`

	// Replicas indicates the number of the pod replicas that will be created.
	// +optional
	// +default=1
	Replicas int32 `json:"replicas,omitempty"`

	// Affinity specifies the scheduling constraints for Pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
}

// RegistryStatus defines the observed state of Registry.
type RegistryStatus struct {
	// Ready is a boolean field that is true when the Registry is ready to be used.
	// +optional
	Ready bool `json:"ready"`

	// Version of the managed Registry.
	// +optional
	Version string `json:"version,omitempty"`

	// Image indicates the container image to use for the Registry.
	// +optional
	Image string `json:"image,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".status.image"

// Registry is the Schema for the registries API.
type Registry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistrySpec   `json:"spec,omitempty"`
	Status RegistryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RegistryList contains a list of Registry.
type RegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Registry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Registry{}, &RegistryList{})
}
