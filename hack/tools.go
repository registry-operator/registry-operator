//go:build never
// +build never

package hack

import (
	_ "github.com/elastic/crd-ref-docs"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/kyverno/chainsaw"
	_ "github.com/tilt-dev/ctlptl/cmd/ctlptl"
	_ "golang.stackrox.io/kube-linter/cmd/kube-linter"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	_ "sigs.k8s.io/kind"
	_ "sigs.k8s.io/kustomize/kustomize/v5"
)
