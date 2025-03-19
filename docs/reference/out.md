# API Reference

## Packages
- [registry-operator.dev/v1alpha1](#registry-operatordevv1alpha1)


## registry-operator.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the registry v1alpha1 API group

### Resource Types
- [Registry](#registry)



#### Registry



Registry is the Schema for the registries API.





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `registry-operator.dev/v1alpha1` | | |
| `kind` _string_ | `Registry` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[RegistrySpec](#registryspec)_ |  |  |  |
| `status` _[RegistryStatus](#registrystatus)_ |  |  |  |


#### RegistrySpec



RegistrySpec defines the desired state of Registry.



_Appears in:_
- [Registry](#registry)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image indicates the container image to use for the Registry. |  |  |
| `replicas` _integer_ | Replicas indicates the number of the pod replicas that will be created. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#resourcerequirements-v1-core)_ | Resources describe the compute resource requirements. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#affinity-v1-core)_ | Affinity specifies the scheduling constraints for Pods. |  |  |
| `storage` _[Storage](#storage)_ | Storage defines the available storage options for a registry.<br />It allows specifying different volume sources to manage storage lifecycle and persistence. |  |  |


#### RegistryStatus



RegistryStatus defines the observed state of Registry.



_Appears in:_
- [Registry](#registry)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready is a boolean field that is true when the Registry is ready to be used. |  |  |
| `version` _string_ | Version of the managed Registry. |  |  |
| `image` _string_ | Image indicates the container image to use for the Registry. |  |  |


#### Storage



Storage specifies various types of volume sources that a registry can use for storage.



_Appears in:_
- [RegistrySpec](#registryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `emptyDir` _[EmptyDirVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#emptydirvolumesource-v1-core)_ | EmptyDir represents a temporary directory that shares a pod's lifetime. |  |  |
| `ephemeral` _[EphemeralVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#ephemeralvolumesource-v1-core)_ | Ephemeral represents a volume that is handled by a cluster storage driver.<br />The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,<br />and deleted when the pod is removed. |  |  |
| `hostPath` _[HostPathVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#hostpathvolumesource-v1-core)_ | HostPath represents a directory on the host. |  |  |
| `persistentVolumeClaim` _[PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaimvolumesource-v1-core)_ | PersistentVolumeClaim represents a reference to a PersistentVolumeClaim in the same namespace. |  |  |
| `persistentVolumeClaimTemplate` _[PersistentVolumeClaimSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaimspec-v1-core)_ | PersistentVolumeClaimTemplate allows creating PVCs dynamically.<br />This defines a PVC template that will be instantiated for the pod. |  |  |


