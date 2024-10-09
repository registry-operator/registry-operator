---
title: Introduction
weight: 0
---

## Understanding Operators

Operators are a method of packaging, deploying, and managing a Kubernetes application. They automate operational tasks, like scaling, backups, updates, and more, through custom controllers that extend Kubernetes functionality. In simpler terms, an Operator watches over your resources and ensures your application maintains a desired state by automating routine tasks.

## Installing the Operator

To install the Registry Operator, run the following commands. This will ensure you're always pulling the latest stable release from the operator’s GitHub repository.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/registry-operator/registry-operator/releases/latest' | jq -r '.tag_name')"
kubectl apply -k "https://github.com/registry-operator/registry-operator/?ref=${LATEST}"
```

This command:

1. Fetches the latest release tag using the GitHub API.
1. Applies the corresponding version of the Registry Operator to your Kubernetes cluster using `kubectl`.

Once installed, the operator will begin monitoring the appropriate resources in your cluster based on the CRDs defined.
