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

package manifestutils

import (
	registryv1alpha1 "github.com/registry-operator/registry-operator/api/v1alpha1"
)

// Annotations return the annotations for OpenTelemetryCollector resources.
func Annotations(instance registryv1alpha1.Registry, filterAnnotations []string) (map[string]string, error) {
	// new map every time, so that we don't touch the instance's annotations
	annotations := map[string]string{}

	if nil != instance.ObjectMeta.Annotations {
		for k, v := range instance.ObjectMeta.Annotations {
			if !IsFilteredSet(k, filterAnnotations) {
				annotations[k] = v
			}
		}
	}

	return annotations, nil
}

// PodAnnotations return the spec annotations for OpenTelemetryCollector pod.
func PodAnnotations(instance registryv1alpha1.Registry, filterAnnotations []string) (map[string]string, error) {
	// new map every time, so that we don't touch the instance's annotations
	podAnnotations := map[string]string{}

	annotations, err := Annotations(instance, filterAnnotations)
	if err != nil {
		return nil, err
	}
	// propagating annotations from metadata.annotations
	for kMeta, vMeta := range annotations {
		if _, found := podAnnotations[kMeta]; !found {
			podAnnotations[kMeta] = vMeta
		}
	}

	return podAnnotations, nil
}
