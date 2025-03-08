// Copyright 2025 The Registry Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Additional copyrights:
// Copyright The OpenTelemetry Authors

package registry

import (
	"context"
	"fmt"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/manifests"
	registryupgrade "github.com/registry-operator/registry-operator/internal/upgrade/registry"
	"github.com/registry-operator/registry-operator/internal/version"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	eventTypeNormal  = "Normal"
	eventTypeWarning = "Warning"

	reasonError         = "Error"
	reasonStatusFailure = "StatusFailure"
	reasonInfo          = "Info"
)

// HandleReconcileStatus handles updating the status of the CRDs managed by the operator.
func HandleReconcileStatus(
	ctx context.Context,
	params manifests.Params,
	registry registryv1alpha1.Registry,
	err error,
) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(2).Info("Updating registry status")
	if err != nil {
		params.Recorder.Event(&registry, eventTypeWarning, reasonError, err.Error())
		return ctrl.Result{}, err
	}
	changed := registry.DeepCopy()

	up := &registryupgrade.VersionUpgrade{
		Version:  version.Get(),
		Client:   params.Client,
		Recorder: params.Recorder,
	}

	upgraded, upgradeErr := up.ManagedInstance(ctx, *changed)
	if upgradeErr != nil {
		// don't fail to allow setting the status
		log.V(2).Error(upgradeErr, "Failed to upgrade the Registry CR")
	}

	changed = &upgraded
	statusErr := UpdateRegistryStatus(ctx, params.Client, changed)
	if statusErr != nil {
		params.Recorder.Event(changed, eventTypeWarning, reasonStatusFailure, statusErr.Error())
		return ctrl.Result{}, statusErr
	}

	statusPatch := client.MergeFrom(&registry)
	if err := params.Client.Status().Patch(ctx, changed, statusPatch); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to apply status changes to the Registry CR: %w", err)
	}

	params.Recorder.Event(changed, eventTypeNormal, reasonInfo, "pplied status changes")
	return ctrl.Result{}, nil
}
