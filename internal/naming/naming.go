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

// Package naming is for determining the names for components (containers, services, ...).
package naming

func ConfigMap(registry, hash string) string {
	return DNSName(Truncate("%s-%s", 63, registry, hash))
}

func ConfigVolume() string {
	return "config"
}

func DistributionConfig() string {
	return "config.yaml"
}

// Container returns the name to use for the container in the pod.
func Container() string {
	return "distribution"
}

// Registry builds the registry (deployment/daemonset) name based on the instance.
func Registry(registry string) string {
	return DNSName(Truncate("%s-registry", 63, registry))
}

// Service builds the service name based on the instance.
func Service(registry string) string {
	return DNSName(Truncate("%s-registry", 63, registry))
}

// RegistryDistributionPort builds the name for default distribution container port.
func RegistryDistributionPort() string {
	return "distribution"
}

// ServiceAccount builds the service account name based on the instance.
func ServiceAccount(registry string) string {
	return DNSName(Truncate("%s-registry", 63, registry))
}
