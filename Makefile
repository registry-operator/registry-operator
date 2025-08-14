# Image URL to use all building/pushing image targets
IMG ?= localhost:5005/registry-operator:dev-$(shell git describe --match='' --always --abbrev=6 --dirty)
PLATFORM ?= linux/$(shell go env GOARCH)
CHAINSAW_ARGS ?=

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: manifests generate api-docs

.PHONY: clean
clean:
	-rm -r bin/

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: test
test: manifests generate ## Run tests.
	go test -coverprofile cover.out ./...

.PHONY: test-e2e
test-e2e: chainsaw ## Run the e2e tests against a k8s instance using Kyverno Chainsaw.
	$(CHAINSAW) test ${CHAINSAW_ARGS}

.PHONY: fmt
fmt: golangci-lint ## Run golangci-lint formatters.
	$(GOLANGCI_LINT) fmt

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linters.
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes.
	$(GOLANGCI_LINT) run --fix

.PHONY: lint-manifests
lint-manifests: kustomize kube-linter ## Run kube-linter on Kubernetes manifests.
	$(KUSTOMIZE) build config/default |\
		$(KUBE_LINTER) lint --config=./config/.kube-linter.yaml -

.PHONY: verify-licenses
verify-licenses: addlicense ## Run addlicense to verify if files have license headers.
	find -type f -name "*.go" ! -path "*/vendor/*" | xargs $(ADDLICENSE) -check || (echo 'Run "make update"' && exit 1)

.PHONY: add-licenses
add-licenses: addlicense ## Run addlicense to append license headers to files missing one.
	find -type f -name "*.go" ! -path "*/vendor/*" | xargs $(ADDLICENSE) -c "The Registry Operator Authors"

.PHONY: hadolint
hadolint: ## Run hadolint on Dockerfile
	$(CONTAINER_TOOL) run --rm -i hadolint/hadolint < Dockerfile

.PHONY: serve-docs
serve-docs: ## Serve dev documentation on port 8000
	$(CONTAINER_TOOL) run \
		--rm \
		--volume $(shell pwd):/app:ro \
		--publish 8000:8000 \
		--pull=always \
		docker.io/library/python:latest \
			bash -c "cd /app && \
				pip install -r docs/requirements.txt && \
				mkdocs serve \
					--dev-addr 0.0.0.0:8000 \
					--livereload"

##@ CI

.PHONY: diff
diff: ## Run git diff-index to check if any changes are made.
	git --no-pager diff HEAD --

VERSION ?= main

.PHONY: publish
publish: ## Runs the script that publishes the latest documentation.
	go run ./hack/cmd/publish -version $(VERSION)

.PHONY: release
release: ## Runs the script that generates new release.
	go run ./hack/cmd/release -version $(VERSION)

##@ Build

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build \
		--platform=${PLATFORM} \
		--tag=${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

.PHONY: build-installer
build-installer: manifests generate kustomize ## Generate a consolidated YAML with CRDs and deployment.
	mkdir -p dist
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default > dist/registry-operator.yaml

##@ Documentation

.PHONY: api-docs
api-docs: crd-ref-docs ## Generate API Reference documentation.
	$(CRD_REF_DOCS) \
		--config=./docs/.crd-ref-docs.yaml \
		--source-path=./api/ \
		--renderer=markdown \
		--output-path=./docs/reference/

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: cluster
cluster: kind ctlptl
	PATH=$(LOCALBIN):$$PATH $(CTLPTL) apply -f hack/kind.yaml
	$(KUBECTL) apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml

.PHONY: cluster-reset
cluster-reset: kind ctlptl
	PATH=$(LOCALBIN):$$PATH $(CTLPTL) delete -f hack/kind.yaml

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: kustomize ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
ADDLICENSE     ?= $(LOCALBIN)/addlicense
CHAINSAW       ?= $(LOCALBIN)/chainsaw
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CRD_REF_DOCS   ?= $(LOCALBIN)/crd-ref-docs
CTLPTL         ?= $(LOCALBIN)/ctlptl
KIND           ?= $(LOCALBIN)/kind
KUBE_LINTER    ?= $(LOCALBIN)/kube-linter
KUSTOMIZE      ?= $(LOCALBIN)/kustomize
GOLANGCI_LINT  ?= $(LOCALBIN)/golangci-lint

## Tool Versions
# renovate: datasource=github-tags depName=google/addlicense
ADDLICENSE_VERSION ?= v1.1.1

# renovate: datasource=github-tags depName=kyverno/chainsaw
CHAINSAW_VERSION ?= v0.2.12

# renovate: datasource=github-tags depName=cert-manager/cert-manager
CERT_MANAGER_VERSION ?= v1.18.2

# renovate: datasource=github-tags depName=kubernetes-sigs/controller-tools
CONTROLLER_TOOLS_VERSION ?= v0.18.0

# renovate: datasource=github-tags depName=elastic/crd-ref-docs
CRD_REF_DOCS_VERSION ?= v0.2.0

# renovate: datasource=github-tags depName=tilt-dev/ctlptl
CTLPTL_VERSION ?= v0.8.42

# renovate: datasource=github-tags depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.4.0

# renovate: datasource=github-tags depName=kubernetes-sigs/kind
KIND_VERSION ?= v0.29.0

# renovate: datasource=github-tags depName=stackrox/kube-linter
KUBE_LINTER_VERSION ?= v0.7.5

# renovate: datasource=github-tags depName=kubernetes-sigs/kustomize
KUSTOMIZE_VERSION ?= v5.6.0

.PHONY: addlicense
addlicense: $(ADDLICENSE)-$(ADDLICENSE_VERSION) ## Download addlicense locally if necessary.
$(ADDLICENSE)-$(ADDLICENSE_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(ADDLICENSE),github.com/google/addlicense,$(ADDLICENSE_VERSION))

.PHONY: chainsaw
chainsaw: $(CHAINSAW)-$(CHAINSAW_VERSION) ## Download chainsaw locally if necessary.
$(CHAINSAW)-$(CHAINSAW_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(CHAINSAW),github.com/kyverno/chainsaw,$(CHAINSAW_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)-$(CONTROLLER_TOOLS_VERSION) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN)-$(CONTROLLER_TOOLS_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: crd-ref-docs
crd-ref-docs: $(CRD_REF_DOCS)-$(CRD_REF_DOCS_VERSION) ## Download crd-ref-docs locally if necessary.
$(CRD_REF_DOCS)-$(CRD_REF_DOCS_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(CRD_REF_DOCS),github.com/elastic/crd-ref-docs,$(CRD_REF_DOCS_VERSION))

.PHONY: ctlptl
ctlptl: $(CTLPTL)-$(CTLPTL_VERSION) ## Download ctlptl locally if necessary.
$(CTLPTL)-$(CTLPTL_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(CTLPTL),github.com/tilt-dev/ctlptl/cmd/ctlptl,$(CTLPTL_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION): $(LOCALBIN)
	./hack/install-golangci-lint.sh $(LOCALBIN) $(GOLANGCI_LINT) $(GOLANGCI_LINT_VERSION)

.PHONY: kind
kind: $(KIND)-$(KIND_VERSION) ## Download kind locally if necessary.
$(KIND)-$(KIND_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(KIND),sigs.k8s.io/kind,$(KIND_VERSION))

.PHONY: kube-linter
kube-linter: $(KUBE_LINTER)-$(KUBE_LINTER_VERSION)
$(KUBE_LINTER)-$(KUBE_LINTER_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(KUBE_LINTER),golang.stackrox.io/kube-linter/cmd/kube-linter,$(KUBE_LINTER_VERSION))

.PHONY: kustomize
kustomize: $(KUSTOMIZE)-$(KUSTOMIZE_VERSION) ## Download kustomize locally if necessary.
$(KUSTOMIZE)-$(KUSTOMIZE_VERSION): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
