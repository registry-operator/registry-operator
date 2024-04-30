//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/kyverno/chainsaw"
	_ "github.com/operator-framework/operator-registry/cmd/opm"
	_ "github.com/operator-framework/operator-sdk/cmd/operator-sdk"
	_ "k8s.io/kubelet"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	_ "sigs.k8s.io/kustomize/kustomize/v5"
)
