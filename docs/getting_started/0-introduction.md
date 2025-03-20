---
title: Introduction
weight: 0
---

## Understanding Operators

Operators are a method of packaging, deploying, and managing a Kubernetes application. They automate operational tasks, like scaling, backups, updates, and more, through custom controllers that extend Kubernetes functionality. In simpler terms, an Operator watches over your resources and ensures your application maintains a desired state by automating routine tasks.

## Installing the Operator

### Using kustomization

This method applies the latest configuration by fetching the latest release tag from GitHub.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/registry-operator/registry-operator/releases/latest' | jq -r '.tag_name')"
kubectl apply -k "https://github.com/registry-operator/registry-operator//config/default?ref=${LATEST}"
```

### Using release manifests

Alternatively, you can deploy the operator using the release manifest directly from GitHub.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/registry-operator/registry-operator/releases/latest' | jq -r '.tag_name')"
kubectl apply -f "https://github.com/registry-operator/registry-operator/releases/download/${LATEST}/registry-operator.yaml"
```

## Updating the Operator

To update to the latest version, rerun the installation command for your chosen method.
