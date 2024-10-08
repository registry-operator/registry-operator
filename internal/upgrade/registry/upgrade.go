// Copyright Registry Operator contributors
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

// Package upgrade handles the upgrade routine from one Registry to the next.
package registry

import (
	"context"
	"fmt"
	"reflect"

	semver "github.com/Masterminds/semver/v3"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
	"github.com/registry-operator/registry-operator/internal/version"

	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type VersionUpgrade struct {
	Client   client.Client
	Recorder record.EventRecorder
	Version  version.Version
}

const RecordBufferSize int = 100

// ManagedInstances finds all the registry instances for the current operator and upgrades them, if necessary.
func (u VersionUpgrade) ManagedInstances(ctx context.Context) error {
	log := log.FromContext(ctx)

	log.Info("Looking for managed instances to upgrade")
	list := &registryv1alpha1.RegistryList{}
	if err := u.Client.List(ctx, list); err != nil {
		return fmt.Errorf("failed to list: %w", err)
	}

	for i := range list.Items {
		original := list.Items[i]
		itemLogger := log.WithValues("registry", klog.KObj(&original))

		upgraded, err := u.ManagedInstance(ctx, original)
		if err != nil {
			const msg = "Automated update not possible. " +
				"Configuration must be corrected manually and CR instance must be re-created."
			itemLogger.Info(msg)
			u.Recorder.Event(&original, "Error", "Upgrade", msg)
			continue
		}
		if !reflect.DeepEqual(upgraded, list.Items[i]) {
			// the resource update overrides the status, so, keep it so that we can reset it later
			st := upgraded.Status
			patch := client.MergeFrom(&original)
			if err := u.Client.Patch(ctx, &upgraded, patch); err != nil {
				itemLogger.Error(err, "failed to apply changes to instance")
				continue
			}

			// the status object requires its own update
			upgraded.Status = st
			if err := u.Client.Status().Patch(ctx, &upgraded, patch); err != nil {
				itemLogger.Error(err, "failed to apply changes to instance's status object")
				continue
			}
			itemLogger.Info("Instance upgraded", "version", upgraded.Status.Version)
		}
	}

	if len(list.Items) == 0 {
		log.Info("No instances to upgrade")
	}

	return nil
}

// ManagedInstance performs the necessary changes to bring the given registry instance to the current version.
func (u VersionUpgrade) ManagedInstance(
	ctx context.Context,
	registry registryv1alpha1.Registry,
) (registryv1alpha1.Registry, error) {
	log := log.FromContext(ctx)

	// this is likely a new instance, assume it's already up to date
	if registry.Status.Version == "" {
		return registry, nil
	}

	instanceV, err := semver.NewVersion(registry.Status.Version)
	if err != nil {
		log.Error(
			err,
			"Failed to parse version for Registry instance",
			"registry", klog.KObj(&registry),
			"version", registry.Status.Version,
		)
		return registry, fmt.Errorf("%w: current version: %v", err, registry.Status.Version)
	}

	updated := *(registry.DeepCopy())
	if instanceV.GreaterThan(&Latest.Version) {
		log.V(4).Info(
			"No upgrade routines are needed for the OpenTelemetry instance",
			"registry", klog.KObj(&updated),
			"version", updated.Status.Version,
			"latest", Latest.Version.String(),
		)

		registryV, err := semver.NewVersion(u.Version.Registry)
		if err != nil {
			return updated, fmt.Errorf("%w: new version: %v", err, u.Version.Registry)
		}

		if instanceV.LessThan(registryV) {
			log.Info(
				"Upgraded Registry version",
				"registry", klog.KObj(&updated),
				"version", updated.Status.Version,
			)
			updated.Status.Version = u.Version.Registry
		} else {
			log.V(4).Info(
				"Skipping upgrade for Registry instance",
				"registry", klog.KObj(&updated),
			)
		}

		return updated, nil
	}

	for _, available := range versions {
		if available.GreaterThan(instanceV) {
			upgraded, err := available.upgrade(u, &updated)
			if err != nil {
				log.Error(
					err,
					"Failed to upgrade managed registry instances",
					"registry", klog.KObj(&updated),
				)
				return updated, err
			}

			log.V(1).Info(
				"Step upgrade",
				"registry", klog.KObj(&updated),
				"version", available.String(),
			)
			upgraded.Status.Version = available.String()
			updated = *upgraded
		}
	}
	// Update with the latest known version, which is what we have from values.yaml
	updated.Status.Version = u.Version.Registry

	log.V(1).Info("Final version", "name", updated.Name, "namespace", updated.Namespace, "version", updated.Status.Version)
	return updated, nil
}
