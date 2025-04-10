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

	// Resources describe the compute resource requirements.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Affinity specifies the scheduling constraints for Pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Storage defines the available storage options for a registry.
	// It allows specifying different storage sources to manage storage lifecycle and persistence.
	// +optional
	Storage Storage `json:"storage,omitempty"`
}

// Storage specifies various types of storage sources that a registry can use for persistence.
type Storage struct {
	// EmptyDir represents a temporary directory that shares a pod's lifetime.
	// +optional
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`

	// Ephemeral represents a volume that is handled by a cluster storage driver.
	// The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,
	// and deleted when the pod is removed.
	// +optional
	Ephemeral *corev1.EphemeralVolumeSource `json:"ephemeral,omitempty"`

	// HostPath represents a directory on the host.
	// +optional
	HostPath *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`

	// PersistentVolumeClaim represents a reference to a PersistentVolumeClaim in the same namespace.
	// +optional
	PersistentVolumeClaim *corev1.PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`

	// PersistentVolumeClaimTemplate allows creating PVCs dynamically.
	// This defines a PVC template that will be instantiated for the pod.
	// +optional
	PersistentVolumeClaimTemplate *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaimTemplate,omitempty"`

	// S3 defines an S3-compatible storage source for persisting registry data.
	// It provides a way to use object storage systems such as Amazon S3 or S3-compatible services
	// for data persistence. This field is optional and can be configured with an endpoint and appropriate credentials.
	// +optional
	S3 *S3StorageSource `json:"s3,omitempty"`
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

// S3StorageSource defines the configuration for connecting to an S3-compatible
// storage backend. It holds the necessary secret references to access an S3-compatible
// storage service.
type S3StorageSource struct {
	// BucketName is an optional reference to the secret key containing the
	// default bucket name to be used.
	BucketName corev1.SecretKeySelector `json:"bucketName"`

	// Region is an optional reference to the secret key containing the S3
	// region name.
	Region corev1.SecretKeySelector `json:"region"`

	// AccessKey is a reference to the secret key containing the S3 access key.
	// +optional
	AccessKey *SecretKeySelector `json:"accessKey,omitempty"`

	// SecretKey is a reference to the secret key containing the S3 secret key.
	// +optional
	SecretKey *SecretKeySelector `json:"secretKey,omitempty"`

	// EndpointURL is an optional reference to the secret key containing an
	// override for the S3 endpoint URL.
	// +optional
	EndpointURL *SecretKeySelector `json:"endpointURL,omitempty"`
}

// SecretKeySelector selects a key of a Secret.
type SecretKeySelector struct {
	// The name of the secret in the object's namespace to select from.
	corev1.LocalObjectReference `json:",inline"`

	// The key of the secret to select from. Must be a valid secret key.
	Key string `json:"key"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:conversion:hub
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
