# API Reference

## Packages
- [registry.registry-operator.dev/v1alpha1](#registryregistry-operatordevv1alpha1)


## registry.registry-operator.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the registry v1alpha1 API group

### Resource Types
- [Registry](#registry)



#### Registry



Registry is the Schema for the registries API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `registry.registry-operator.dev/v1alpha1` | | |
| `kind` _string_ | `Registry` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/vv1.31.1/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[RegistrySpec](#registryspec)_ |  |  |  |
| `status` _[RegistryStatus](#registrystatus)_ |  |  |  |


#### RegistrySpec



RegistrySpec defines the desired state of Registry



_Appears in:_
- [Registry](#registry)



#### RegistryStatus



RegistryStatus defines the observed state of Registry



_Appears in:_
- [Registry](#registry)



