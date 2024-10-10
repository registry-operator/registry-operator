// Copyright The Registry Operator Authors
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
	"github.com/Masterminds/semver/v3"

	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
)

type upgradeFunc func(u VersionUpgrade, registry *registryv1alpha1.Registry) (*registryv1alpha1.Registry, error)

type registryVersion struct {
	upgrade upgradeFunc
	semver.Version
}

var (
	versions = []registryVersion{
		{
			Version: *semver.MustParse("2.8.3"),
			upgrade: upgrade2_8_3,
		},
	}

	// Latest represents the latest version that we need to upgrade. This is not necessarily the latest known version.
	Latest = versions[len(versions)-1]
)
