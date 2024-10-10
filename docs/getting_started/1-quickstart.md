---
title: Quick Start
weight: 1
---

## Deploying the simplest Registry

The simplest possible way to create a Registry instance is by creating a YAML file like the following example. This will install the Registry using default latest[^1] image in a **single pod**, using **in-memory storage** by default, with a [`ClusterIP`][k8s/svc/cluster-ip] service.

```yaml
apiVersion: registry-operator.dev/v1alpha1
kind: Registry
metadata:
  name: simplest
```

The YAML file can then be used with `kubectl`:

```sh
kubectl apply -f simplest.yaml
```

In a few seconds, a new in-memory instance of Registry will be available, suitable for quick demos and development purposes. To check the instances that were created, list the Registry objects:

```
$ kubectl get registries.registry-operator.dev 
NAME       VERSION   READY   IMAGE
simplest   2.8.3     true    docker.io/library/registry:2.8.3
```

To get the deplpyment name, query for the deployments belonging to the simplest Registry instance:

```
$ kubectl get deployments.apps -l=app.kubernetes.io/instance=default.simplest
NAME                READY   UP-TO-DATE   AVAILABLE   AGE
simplest-registry   1/1     1            1           2m
```

To get the service name, query for the services belonging to the simplest Registry instance:

```
$ kubectl get services -l=app.kubernetes.io/instance=default.simplest
NAME                TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)    AGE
simplest-registry   ClusterIP   10.96.187.2   <none>        5000/TCP   2m
```

Similarly, the logs can be queried from all pods belonging to our instance:

```
$ kubectl logs -l app.kubernetes.io/instance=default.simplest
time="2024-10-09T08:11:02.799689427Z" level=warning msg="No HTTP secret provided - generated random secret. This may cause problems with uploads if multiple registries are behind a load-balancer. To provide a shared secret, fill in http.secret in the configuration file or set the REGISTRY_HTTP_SECRET environment variable." go.version=go1.20.8 instance.id=8ab1846a-75f9-4ac0-840e-4876165f56b2 service=registry version=2.8.3 
time="2024-10-09T08:11:02.799779511Z" level=info msg="redis not configured" go.version=go1.20.8 instance.id=8ab1846a-75f9-4ac0-840e-4876165f56b2 service=registry version=2.8.3 
time="2024-10-09T08:11:02.800089011Z" level=info msg="using inmemory blob descriptor cache" go.version=go1.20.8 instance.id=8ab1846a-75f9-4ac0-840e-4876165f56b2 service=registry version=2.8.3 
time="2024-10-09T08:11:02.800370261Z" level=info msg="listening on [::]:5000" go.version=go1.20.8 instance.id=8ab1846a-75f9-4ac0-840e-4876165f56b2 service=registry version=2.8.3 
time="2024-10-09T08:11:02.800428594Z" level=info msg="Starting upload purge in 16m0s" go.version=go1.20.8 instance.id=8ab1846a-75f9-4ac0-840e-4876165f56b2 service=registry version=2.8.3 
```

[^1]: Latest image is the latest stable version available at the time of the Registry Operator release.
[k8s/svc/cluster-ip]: https://kubernetes.io/docs/concepts/services-networking/service/#type-clusterip
