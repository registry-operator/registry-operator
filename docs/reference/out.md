# API Reference

## Packages
- [registry-operator.dev/v1alpha1](#registry-operatordevv1alpha1)


## registry-operator.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the registry v1alpha1 API group

### Resource Types
- [Registry](#registry)



#### Auth



Auth



_Appears in:_
- [RegistrySpec](#registryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `htpasswd` _[Htpasswd](#htpasswd)_ | Htpasswd |  |  |


#### Htpasswd



Htpasswd



_Appears in:_
- [Auth](#auth)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `realm` _string_ | Realm |  |  |
| `secret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#secretkeyselector-v1-core)_ | Htpasswd |  |  |


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
| `auth` _[Auth](#auth)_ | Auth |  |  |
| `storage` _[Storage](#storage)_ | Storage defines the available storage options for a registry.<br />It allows specifying different storage sources to manage storage lifecycle and persistence. |  |  |


#### RegistryStatus



RegistryStatus defines the observed state of Registry.



_Appears in:_
- [Registry](#registry)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ready` _boolean_ | Ready is a boolean field that is true when the Registry is ready to be used. |  |  |
| `version` _string_ | Version of the managed Registry. |  |  |
| `image` _string_ | Image indicates the container image to use for the Registry. |  |  |


#### S3StorageSource



S3StorageSource defines the configuration for connecting to an S3-compatible
storage backend. It holds the necessary secret references to access an S3-compatible
storage service.



_Appears in:_
- [Storage](#storage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `bucketName` _[SecretKeySelector](#secretkeyselector)_ | BucketName is an optional reference to the secret key containing the<br />default bucket name to be used. |  |  |
| `region` _[SecretKeySelector](#secretkeyselector)_ | Region is an optional reference to the secret key containing the S3<br />region name. |  |  |
| `accessKey` _[SecretKeySelector](#secretkeyselector)_ | AccessKey is a reference to the secret key containing the S3 access key. |  |  |
| `secretKey` _[SecretKeySelector](#secretkeyselector)_ | SecretKey is a reference to the secret key containing the S3 secret key. |  |  |
| `endpointURL` _[SecretKeySelector](#secretkeyselector)_ | EndpointURL is an optional reference to the secret key containing an<br />override for the S3 endpoint URL. |  |  |


#### SecretKeySelector



SecretKeySelector selects a key of a Secret.



_Appears in:_
- [S3StorageSource](#s3storagesource)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `key` _string_ | The key of the secret to select from. Must be a valid secret key. |  |  |


#### Storage



Storage specifies various types of storage sources that a registry can use for persistence.



_Appears in:_
- [RegistrySpec](#registryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `emptyDir` _[EmptyDirVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#emptydirvolumesource-v1-core)_ | EmptyDir represents a temporary directory that shares a pod's lifetime. |  |  |
| `ephemeral` _[EphemeralVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#ephemeralvolumesource-v1-core)_ | Ephemeral represents a volume that is handled by a cluster storage driver.<br />The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,<br />and deleted when the pod is removed. |  |  |
| `hostPath` _[HostPathVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#hostpathvolumesource-v1-core)_ | HostPath represents a directory on the host. |  |  |
| `persistentVolumeClaim` _[PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaimvolumesource-v1-core)_ | PersistentVolumeClaim represents a reference to a PersistentVolumeClaim in the same namespace. |  |  |
| `persistentVolumeClaimTemplate` _[PersistentVolumeClaimSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaimspec-v1-core)_ | PersistentVolumeClaimTemplate allows creating PVCs dynamically.<br />This defines a PVC template that will be instantiated for the pod. |  |  |
| `s3` _[S3StorageSource](#s3storagesource)_ | S3 defines an S3-compatible storage source for persisting registry data.<br />It provides a way to use object storage systems such as Amazon S3 or S3-compatible services<br />for data persistence. This field is optional and can be configured with an endpoint and appropriate credentials. |  |  |


